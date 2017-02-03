package groupcache

import (
	"bytes"
	"demo/groupcache/consistenthash"
	"demo/groupcache/groupcachepb"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/golang/protobuf/proto"
)

const defaultBasePath = "/_groupcache/"
const defaultReplicas = 50

type HTTPPool struct {
	//当请求到来时为服务器指定一个上下文,
	Context func(*http.Request) Context
	//创建请求时为客户端设定一个http.RoundTripper
	Transport func(Context) http.RoundTripper
	//perr的基础URL,例如 https://example.com:8888
	self string
	//选项
	opts HTTPPoolOptions

	mu          sync.Mutex //保护peers和httpGetters
	peers       *consistenthash.Map
	httpGetters map[string]*httpGetter //key是URL,例如 http://10.0.0.2:8888
}

type HTTPPoolOptions struct {
	BasePath string              //服务groupcache请求的HTTP path
	Replicas int                 //哈希计算时使用
	HashFn   consistenthash.Hash //哈希函数,默认crc32.ChecksumIEEE
}

func NewHTTPPool(self string) *HTTPPool {
	//生成一个PeerPicker
	p := NewHTTPPoolOpts(self, nil)
	//进行Handle注册
	http.Handle(p.opts.BasePath, p)
	return p
}

var (
	httpPoolMade bool
)

func NewHTTPPoolOpts(self string, o *HTTPPoolOptions) *HTTPPool {
	if httpPoolMade {
		panic("groupcache: NewHTTPPool must be called only once")
	}
	httpPoolMade = true

	p := &HTTPPool{
		self:        self,
		httpGetters: make(map[string]*httpGetter),
	}

	if o != nil {
		p.opts = *o
	}
	if p.opts.BasePath == "" {
		p.opts.BasePath = defaultBasePath
	}
	if p.opts.Replicas == 0 {
		p.opts.Replicas = defaultReplicas
	}
	//hash map 存储其它peer
	p.peers = consistenthash.New(p.opts.Replicas, p.opts.HashFn)
	//此处注册peer查找函数,直接返回当前的HTTPPool
	RegisterPeerPicker(func() PeerPicker { return p })
	//	display.Display("NEW_HTTPPool", p)
	return p
}

//更新pool的基础url列表,重新设置peers
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.New(p.opts.Replicas, p.opts.HashFn)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{transport: p.Transport, baseURL: peer + p.opts.BasePath}
	}
}

//取出一个peer
func (p *HTTPPool) PickPeer(key string) (ProtoGetter, bool) {
	//	display.Display("HTTPPool", p)
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.peers.IsEmpty() {
		return nil, false
	}
	//不是自身
	if peer := p.peers.Get(key); peer != p.self {
		logger.Infof("----peer : %s , httpGetter.baseURL : %s\n", peer, p.httpGetters[peer].baseURL)
		return p.httpGetters[peer], true
	}
	return nil, false
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//	display.Display("HTTPPool", p)
	logger.Info(r.URL)
	//解析URL
	if !strings.HasPrefix(r.URL.Path, p.opts.BasePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	//用"/"分割Path,返回3-last http://example.com/hello 返回example.com/hello
	parts := strings.SplitN(r.URL.Path[len(p.opts.BasePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group"+groupName, http.StatusNotFound)
		return
	}
	var ctx Context
	if p.Context != nil {
		ctx = p.Context(r)
	}

	group.Stats.ServerRequests.Add(1)
	var value []byte
	err := group.Get(ctx, key, AllocatingByteSliceSink(&value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := proto.Marshal(&groupcachepb.GetResponse{Value: value})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/x-protobuf")
	w.Write(body)
}

type httpGetter struct {
	transport func(Context) http.RoundTripper
	baseURL   string
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

//peer.Get实际调用函数
func (h *httpGetter) Get(context Context, in *groupcachepb.GetRequest, out *groupcachepb.GetResponse) error {
	//	display.Display("httpGetter", h)
	logger.Infof("enter into %v.Get  http method", h.baseURL)
	defer logger.Infof("leave from %v.Get http method", h.baseURL)
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(in.GetGroup()),
		url.QueryEscape(in.GetKey()),
	)
	// u --> http://127.0.0.1:33820/_groupcache/httpPoolTest/98
	req, err := http.NewRequest("GET", u, nil)
	/*
		Display req (*http.Request):
		(*req).Method = "GET"
		(*(*req).URL).Scheme = "http"
		(*(*req).URL).Opaque = ""
		(*(*req).URL).User = nil
		(*(*req).URL).Host = "127.0.0.1:43643"
		(*(*req).URL).Path = "/_groupcache/httpPoolTest/99"
		(*(*req).URL).RawPath = ""
		(*(*req).URL).ForceQuery = false
		(*(*req).URL).RawQuery = ""
		(*(*req).URL).Fragment = ""
		(*req).Proto = "HTTP/1.1"
		(*req).ProtoMajor = 1
		(*req).ProtoMinor = 1
		(*req).Body = nil
		(*req).ContentLength = 0
		(*req).Close = false
		(*req).Host = "127.0.0.1:43643"
		(*req).MultipartForm = nil
		(*req).RemoteAddr = ""
		(*req).RequestURI = ""
		(*req).TLS = nil
		(*req).Cancel = <-chan struct {}0x0
		(*req).Response = nil
		(*req).ctx = nil
	*/
	if err != nil {
		return nil
	}
	tr := http.DefaultTransport
	//	display.Display("tr", tr)
	//fmt.Println(h.transport)
	if h.transport != nil {
		tr = h.transport(context)
	}
	res, err := tr.RoundTrip(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", res.Status)
	}
	b := bufferPool.Get().(*bytes.Buffer)
	b.Reset()
	defer bufferPool.Put(b)
	_, err = io.Copy(b, res.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}
	err = proto.Unmarshal(b.Bytes(), out)
	// out	value:"2:127.0.0.1:41118:99"
	if err != nil {
		return fmt.Errorf("decoding response body: %v", err)
	}
	return nil
}
