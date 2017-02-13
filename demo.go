/* Demonstrate how to use rest.RouteObjectMethod
rest.RouteObjectMethod helps create a Route that points to
an object method instead of just a function.
The Curl Demo:
        curl -i -d '{"Name":"Antoine"}' http://127.0.0.1:8080/users
        curl -i http://127.0.0.1:8080/users/0
        curl -i -X PUT -d '{"Name":"Antoine Imbert"}' http://127.0.0.1:8080/users/0
        curl -i -X DELETE http://127.0.0.1:8080/users/0
        curl -i http://127.0.0.1:8080/users
*/
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"bytes"
	"io"
	"time"
	"crypto/md5"
	"strconv"
	"net/url"
	"github.com/golang/protobuf/proto"
	"sync"
	"sort"
)

func sayHelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()                       //解析参数,默认不会解析
	fmt.Println(r.Form)                 //输出到服务器端的信息
	fmt.Println("path", r.URL.Path)     //
	fmt.Println("schema", r.URL.Scheme) //
	fmt.Println(r.Form["url_long"])

	for k, v := range r.Form {
		fmt.Println("key : ", k)
		fmt.Println("value : ", strings.Join(v, ""))
	}

	fmt.Fprintf(w, "hello golang")
}

func main() {
	//注册路由处理函数
	http.HandleFunc("/", sayHelloName) //设置访问路由
	http.HandleFunc("/login", login)   //设置访问路由
	http.HandleFunc("/upload", upload)
	//设置监听端口,和handle
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		log.Fatal("ListenAndServe : ", err)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method : ", r.Method) //获取请求方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		r.ParseForm() //解析参数,默认不会解析
		fmt.Println("username : ", r.Form["username"])
		fmt.Println("password : ", r.Form["password"])
		fmt.Println(r.Form)
	} else if r.Method == "PUT" {

	} else if r.Method == "DELETE" {

	}
}

// 处理/upload 逻辑
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else if r.Method == "POST" {
		fmt.Println("---------\t"  + newUrl(r.URL))
/*
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)  // 此处假设当前目录下已存在test目录
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
*/
	} else if r.Method == "PUT" {
		url := r.URL
		nUrl := url.Scheme
		host := r.Host
		host += ""
		nReq, err := http.NewRequest("PUT", nUrl, nil)
		if err != nil {
			panic(err)
		}
		tr := http.DefaultTransport
		res, err := tr.RoundTrip(nReq)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			fmt.Errorf("server returned: %v", res.Status)
		}
		b := bufferPool.Get().(*bytes.Buffer)
		b.Reset()
		defer bufferPool.Put(b)
		_, err = io.Copy(b, res.Body)
		if err != nil {
			fmt.Errorf("reading response body: %v", err)
		}
		w.Write(b.Bytes())
	} else if r.Method == "DELETE" {

	}
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func newUrl(url *url.URL) string {
	nUrl := ""
	if url.Scheme != "" {
		nUrl += url.Scheme
	} else {
		nUrl += "http"
	}
	nUrl += "://"
	host := url.Host
	if host == "" {
		host = "127.0.0.1:19999"
	} else {
		hosts := strings.Split(host, ":")
		if hosts[0] == "" {
			hosts[0] = "127.0.0.1"
		}
		if hosts[1] == "" {
			hosts[1] = "19999"
		}
		host = strings.Join(hosts, ":")
	}
	nUrl += host + url.Path
	if url.RawQuery != "" {
		nUrl +=  "?" + url.RawQuery
	}
	return nUrl
}

func genSign(url url.URL) string {
	querys := url.RawQuery
	queryList := strings.Split(querys, "&")
	queryMap := make(map[string]string, len(queryList))
	for _, v := range queryList {
		kv := strings.Split(v, "=")
		queryMap[kv[0]] = kv[1]
	}
	queryKeys := make([]string, len(queryList))
	for k, _ := range queryMap {
		queryKeys = append(queryKeys, k)
	}
	sort.Strings(queryKeys)
	newQeury := url.Path + "?"
	for _, k := range queryKeys {
		newQeury += k + "=" + queryMap[k]
	}
	m5 := md5.New()
	m5s := m5.Sum([]byte(newQeury))
	return string(m5s)
}
/*
type URL struct {
	Scheme     string
	Opaque     string    // encoded opaque data
	User       *Userinfo // username and password information
	Host       string    // host or host:port
	Path       string
	RawPath    string // encoded path hint (Go 1.5 and later only; see EscapedPath method)
	ForceQuery bool   // append a query ('?') even if RawQuery is empty
	RawQuery   string // encoded query values, without '?'
	Fragment   string // fragment for references, without '#'
}

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