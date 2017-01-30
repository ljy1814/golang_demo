package groupcache

import (
	"demo/display"
	"demo/groupcache/groupcachepb"
	"demo/groupcache/lru"
	"demo/groupcache/singleflight"
	"demo/mylog"
	"errors"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

//通过key获取数据
type Getter interface {
	Get(ctx Context, key string, dest Sink) error
}

type GetterFunc func(ctx Context, key string, dest Sink) error

func (f GetterFunc) Get(ctx Context, key string, dest Sink) error {
	return f(ctx, key, dest)
}

var (
	mu                 sync.RWMutex
	groups             = make(map[string]*Group)
	initPeerServerOnce sync.Once
	initPeerServer     func()
	logger             = mylog.NewFileLogger()
)

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	return newGroup(name, cacheBytes, getter, nil)
}

//新生成一个group
func newGroup(name string, cacheBytes int64, getter Getter, peers PeerPicker) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	initPeerServerOnce.Do(callInitPeerServer)
	if _, dup := groups[name]; dup {
		panic("duplicate registeration of group " + name)
	}
	g := &Group{
		name:       name,
		getter:     getter,
		peers:      peers,
		cacheBytes: cacheBytes,
		loadGroup:  &singleflight.Group{},
	}
	if fn := newGroupHook; fn != nil {
		fn(g)
	}
	groups[name] = g
	//	logger.Infoln(groups)
	GDisplay("groups", groups)
	return g
}

func GDisplay(name string, value interface{}) {
	fdisplay, err := os.OpenFile("display-"+time.Now().Format("2006-01-02")+".log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic("open display filoe failed")
	}
	defer fdisplay.Close()
	display.Fdisplay(fdisplay, name, value)
}

var newGroupHook func(*Group)

func RegisterNewGroupHook(fn func(*Group)) {
	if newGroupHook != nil {
		panic("RegisterNewGroupHook called more than once")
	}
	newGroupHook = fn
}

func RegisterServerStart(fn func()) {
	if initPeerServer != nil {
		panic("RegisterServerStart called more than once")
	}
	initPeerServer = fn
}

func callInitPeerServer() {
	if initPeerServer != nil {
		initPeerServer()
	}
}

type Group struct {
	name       string
	getter     Getter
	peersOnce  sync.Once
	peers      PeerPicker
	cacheBytes int64 //mainCache和hotCache总和为其最大值

	//哈希有效的key集合
	mainCache cache
	//访问频率最高的key/value集合,包含不在mainCache内的
	hotCache cache

	//保证每个key只被取一次
	loadGroup flightGroup

	_     int32 //在32位平台上保证8字节对其
	Stats Stats
}

type flightGroup interface {
	Do(key string, fn func() (interface{}, error)) (interface{}, error)
}

type Stats struct {
	Gets            AtomicInt //任何Get请求,包括来自对等体的
	CacheHits       AtomicInt //任何cache都行
	PeerLoads       AtomicInt //包括远程加载和远程缓存获取
	PeerErrors      AtomicInt
	Loads           AtomicInt //获取缓存
	LoadsDeduped    AtomicInt //在singleflight之后
	LocalLoads      AtomicInt //所有的本地正常加载的
	LocalLoadErrors AtomicInt //加载错误的
	ServerRequests  AtomicInt //网上的请求
}

func (g *Group) Name() string {
	return g.name
}

func (g *Group) initPeers() {
	if g.peers == nil {
		g.peers = getPeers(g.name)
	}
}

func (g *Group) Get(ctx Context, key string, dest Sink) error {
	logger.Infof("enter into %v.Get method", g.name)
	defer logger.Infof("leave from %v.Get method", g.name)
	g.peersOnce.Do(g.initPeers)
	//get计数加1
	g.Stats.Gets.Add(1)
	if dest == nil {
		return errors.New("groupcache: nil dest Sink")
	}
	value, cacheHit := g.lookupCache(key)

	//命中缓存
	if cacheHit {
		g.Stats.CacheHits.Add(1)
		return setSinkView(dest, value)
	}

	//两个cache都没有查找到
	//避免两次反序列化或者复制,多goroutine同时访问
	destPopulated := false
	value, destPopulated, err := g.load(ctx, key, dest)
	if err != nil {
		return err
	}
	if destPopulated {
		return nil
	}
	return setSinkView(dest, value)
}

//加载数据,本地加载或者从网络上获取
func (g *Group) load(ctx Context, key string, dest Sink) (value ByteView, destPopulated bool, err error) {
	logger.Infof("enter into %v.load method", g.name)
	defer logger.Infof("leave from %v.load method", g.name)
	g.Stats.Loads.Add(1)
	viewi, err := g.loadGroup.Do(key, func() (interface{}, error) {
		//查询cache,防止两个goroutine同时访问出现差错,singleflight只会调用那个函数
		/*
			1 Get("key")
			2 Get("key")
			1 lookupCache("key")
			2 lookupCache("key")
			1 load("key")
			2 load("key")
			1 loadGroup("key", fn)
			1 fn
			2 loadGroup("key", fn)
			2 fn
		*/
		//防止有一个goroutine已经将其加载到缓存里面了
		if value, cacheHit := g.lookupCache(key); cacheHit {
			g.Stats.CacheHits.Add(1)
			return value, nil
		}
		g.Stats.LoadsDeduped.Add(1)
		var value ByteView
		var err error
		//获取网络上的对等体,同类客户端或者服务器
		if peer, ok := g.peers.PickPeer(key); ok {
			value, err = g.getFromPeer(ctx, peer, key)
			if err == nil {
				g.Stats.PeerLoads.Add(1)
				return value, nil
			}
			//其它对等体没有查找到,
			g.Stats.PeerErrors.Add(1)
		}
		//没有从网络上查找到,从本地查找
		value, err = g.getLocally(ctx, key, dest)
		if err != nil {
			g.Stats.LocalLoadErrors.Add(1)
			return nil, err
		}
		g.Stats.LocalLoads.Add(1)
		destPopulated = true
		g.populateCache(key, value, &g.mainCache)
		return value, nil
	})
	if err == nil {
		value = viewi.(ByteView)
	}
	return
}

func (g *Group) getLocally(ctx Context, key string, dest Sink) (ByteView, error) {
	logger.Infof("enter into %v.getLocally method", g.name)
	defer logger.Infof("leave from %v.getLocally method", g.name)
	err := g.getter.Get(ctx, key, dest)
	if err != nil {
		return ByteView{}, err
	}
	return dest.view()
}

func (g *Group) getFromPeer(ctx Context, peer ProtoGetter, key string) (ByteView, error) {
	logger.Infof("enter into %v.getFromPeer method", g.name)
	defer logger.Infof("leave from %v.getFromPeer method", g.name)
	req := &groupcachepb.GetRequest{
		Group: &g.name,
		Key:   &key,
	}
	res := &groupcachepb.GetResponse{}
	err := peer.Get(ctx, req, res)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: res.Value}

	//随机清理cache中的元素
	if rand.Intn(10) == 0 {
		g.populateCache(key, value, &g.hotCache)
	}
	return value, nil
}

func (g *Group) lookupCache(key string) (value ByteView, ok bool) {
	logger.Infof("enter into %v.lookupCache method", g.name)
	defer logger.Infof("leave from %v.lookupCache method", g.name)
	if g.cacheBytes <= 0 {
		return
	}
	//先从mainCache里面查找,没有则去hotCache里面查找
	value, ok = g.mainCache.get(key)
	if ok {
		return
	}
	value, ok = g.hotCache.get(key)
	return
}

func (g *Group) populateCache(key string, value ByteView, cache *cache) {
	logger.Infof("enter into %v.popyulateCache method", g.name)
	defer logger.Infof("leave from %v.populateCache method", g.name)
	if g.cacheBytes <= 0 {
		return
	}
	cache.add(key, value)
	//Evict items from cache(s) if necessary
	for {
		mainBytes := g.mainCache.bytes()
		hotBytes := g.hotCache.bytes()
		//缓存未满
		if mainBytes+hotBytes <= g.cacheBytes {
			return
		}
		victim := &g.mainCache
		//到了一定程度就清理最旧的元素
		if hotBytes > mainBytes/8 {
			victim = &g.hotCache
		}
		victim.removeOldest()
	}
}

type CacheType int

const (
	MainCache CacheType = iota + 1
	HotCache
)

func (g *Group) CacheStats(which CacheType) CacheStats {
	switch which {
	case MainCache:
		return g.mainCache.stats()
	case HotCache:
		return g.hotCache.stats()
	default:
		return CacheStats{}
	}
}

type cache struct {
	mu         sync.RWMutex
	nbytes     int64 //所有的键和值
	lru        *lru.Cache
	nhit, nget int64
	nevict     int64
}

func (c *cache) stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return CacheStats{
		Bytes:     c.nbytes,
		Items:     c.itemsLocked(),
		Gets:      c.nget,
		Hits:      c.nhit,
		Evictions: c.nevict,
	}
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = &lru.Cache{
			OnEvicted: func(key lru.Key, value interface{}) {
				val := value.(ByteView)
				c.nbytes -= int64(len(key.(string))) + int64(val.Len())
				c.nevict++
			},
		}
	}
	c.lru.Add(key, value)
	c.nbytes += int64(len(key)) + int64(value.Len())
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.nget++
	if c.lru == nil {
		return
	}
	vi, ok := c.lru.Get(key)
	if !ok {
		return
	}
	c.nhit++
	return vi.(ByteView), true
}

func (c *cache) removeOldest() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru != nil {
		c.lru.RemoveOldest()
	}
}

func (c *cache) bytes() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.nbytes
}

func (c *cache) items() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.itemsLocked()
}

func (c *cache) itemsLocked() int64 {
	if c.lru == nil {
		return 0
	}
	return int64(c.lru.Len())
}

type AtomicInt int64

func (i *AtomicInt) Add(n int64) {
	atomic.AddInt64((*int64)(i), n)
}
func (i *AtomicInt) Get() int64 {
	return atomic.LoadInt64((*int64)(i))
}

func (i *AtomicInt) String() string {
	return strconv.FormatInt(i.Get(), 10)
}

type CacheStats struct {
	Bytes     int64
	Items     int64
	Gets      int64
	Hits      int64
	Evictions int64
}
