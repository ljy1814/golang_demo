/* Demonstrate how to use rest.RouteObjectMethod
rest.RouteObjectMethod helps create a Route that points to
an object method instead of just a function.
The Curl Demo:
        curl -i -d '{"Name":"Antoine"}' http://127.0.0.1:8080/users
        curl -i http://127.0.0.1:8080/users/0
        curl -i -X PUT -d '{"Name":"Antoine Imbert"}' http://127.0.0.1:8080/users/0
        curl -i -X DELETE http://127.0.0.1:8080/users/0
        curl -i http://127.0.0.1:8080/users

URL ----   http://www.example.com/user/info/123/update?type=1&number=111
URL ----   http://www.example.com/user/search?type=1&nums=123&title=haha&age=34&time=1434524242
URL ----   http://www.example.com/user/search?type=1&nums=123&title=haha&age=34&time=1434524242&channel=asadda&sign=12343131
method=GET PUT POST DELETE
contentType=JSON,FILE,FORM
content='{"AAA":"BBB"}'

name=xxx
age=123
*/
package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
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
	//	Display("req", r)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		//		r.ParseForm() //解析参数,默认不会解析
		fmt.Println("username : ", r.Form["username"])
		fmt.Println("password : ", r.Form["password"])
		fmt.Println(r.Form)
		fmt.Println(r.Body)
		nUrl := handleUrl(r)
		fmt.Printf("xxxx nUrl : %s\n", nUrl)
		r.ParseForm()
		fmt.Println(r.Form)
		//		nReq, err := http.NewRequest("POST", nUrl, strings.NewReader(r.Form.Encode()))
		//		nReq, err := http.NewRequest("POST", nUrl, nil)
		//		if err != nil {
		//			panic(err)
		//		}
		//		handleReq(r, nReq)
		//		res, err := http.Post(nUrl, "application/x-www-form-urlencoded", strings.NewReader(r.Form.Encode()))
		//		nReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		//		client := &http.Client{}
		//		res, err := client.Do(nReq)
		//		tr := http.DefaultTransport
		//		res, err := tr.RoundTrip(nReq)
		nReq, err := http.NewRequest("POST", nUrl, strings.NewReader(r.Form.Encode()))
		if err != nil {
			panic(err)
		}
		nReq.Header.Set("contentType", r.Form["contentType"][0])

		client := &http.Client{}
		res, err := client.Do(nReq)
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
		fmt.Println(b)
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

		url := r.URL
		nUrl := newUrl(url)
		if strings.Index(nUrl, "?") != len(nUrl)-1 {
			nUrl += "&type=" + getType(r)
		} else {
			nUrl += "type=" + getType(r)
		}
		nUrl += "&sign=" + genSign(url)
		fmt.Println(nUrl)
		//nReq, err := http.NewRequest("PUT", nUrl, nil)
		//if err != nil {
		//	panic(err)
		//}
		//		handleReq(r, nReq)
		//		tr := http.DefaultTransport
		//		res, err := tr.RoundTrip(nReq)
		r.ParseMultipartForm(32 << 20)
		fmt.Println(r.Form)
		fmt.Println(r.FormFile("uploadfile"))
		fff, mff, err := r.FormFile("uploadfile")
		fmt.Println(fff)
		fmt.Println(mff)
		if err != nil {
			panic(err)
		}
		res := postFile(mff.Filename, nUrl)
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
		/*		res, err := http.Post(nUrl, "application/x-www-form-urlencoded", strings.NewReader(r.Form.Encode()))
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
				fmt.Println(b)
		*/

		/*   server
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
	}
}

func postFile(filename string, targetUrl string) *http.Response {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//TODO important
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		fmt.Println("error write to file")
		panic(err)
	}

	//open file
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		panic(err)
	}

	//copy file to proxy server
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		panic(err)
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	nReq, err := http.NewRequest("POST", targetUrl, bodyBuf)
	nReq.Header.Set("Content-Type", contentType)
	nReq.Header.Set("contentType", "file")

	client := &http.Client{}
	resp, err := client.Do(nReq)
	if err != nil {
		panic(err)
	}
	return resp

	/*
		resp, err := http.Post(targetUrl, contentType, bodyBuf)
		if err != nil {
			return err
		}
	*/

	//TODO temp
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))
	return nil
}

func handleReq(or *http.Request, r *http.Request) {
	r.Form = or.Form
	r.MultipartForm = or.MultipartForm
	r.PostForm = or.PostForm
	r.Header = or.Header
	r.Body = or.Body
	r.Close = or.Close
}

func handleUrl(r *http.Request) string {
	url := r.URL
	nUrl := newUrl(url)
	if strings.Index(nUrl, "?") != len(nUrl)-1 {
		nUrl += "channel=" + fmt.Sprintf("%x", md5.New().Sum(nil))
	} else {
		nUrl += "&channel=" + fmt.Sprintf("%x", md5.New().Sum(nil))
	}
	nUrl += "&sign=" + genSign(url)
	fmt.Printf("nUrl : %s\n", nUrl)
	return nUrl
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

/*
0 FORM
1 JSON
2 FILE
*/
func getType(r *http.Request) string {
	r.ParseForm()
	r.ParseMultipartForm(32 << 20)
	result, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	fmt.Println("result : " + string(result))
	if r.MultipartForm != nil {
		return "file"
	}
	//未知类型的推荐处理方法
	var f interface{}
	//	err := json.NewDecoder(r.Body).Decode(&f)
	for k, _ := range r.Form {
		err := json.Unmarshal([]byte(k), &f)
		fmt.Println(err)
		if err == nil {
			return "json"
		}
		break
	}
	return "form"
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
	nUrl += "?" + url.RawQuery
	return nUrl
}

func genSign(url *url.URL) string {
	querys := url.RawQuery
	newQeury := url.Path + "?"
	if querys != "" {
		queryList := strings.Split(querys, "&")
		queryMap := make(map[string]string, len(queryList))
		for _, v := range queryList {
			if v != "" {
				kv := strings.Split(v, "=")
				queryMap[kv[0]] = kv[1]
			}

		}
		queryKeys := make([]string, len(queryList))
		for k, _ := range queryMap {
			queryKeys = append(queryKeys, k)
		}
		sort.Strings(queryKeys)

		for _, k := range queryKeys {
			newQeury += k + "=" + queryMap[k]
		}
	}
	m5 := md5.New()
	token := fmt.Sprintf("%x", m5.Sum([]byte(newQeury)))
	return token
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

func display(path string, v reflect.Value) {
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Printf("%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			display(fmt.Sprintf("%s[%d]", path, i), v.Index(i))
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			display(fieldPath, v.Field(i))
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			display(fmt.Sprintf("%s[%s]", path, formatAtom(key)), v.MapIndex(key))
		}
	case reflect.Ptr:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			display(fmt.Sprintf("(*%s)", path), v.Elem())
		}
	case reflect.Interface:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			fmt.Printf("%s.type = %s\n", path, v.Elem().Type())
			display(path+".value", v.Elem())
		}
	default:
		fmt.Printf("%s = %s\n", path, formatAtom(v))
	}
}

func Display(name string, x interface{}) {
	fmt.Printf("Display %s (%T):\n", name, x)
	display(name, reflect.ValueOf(x))
}

func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		return v.Type().String() + "0x" + strconv.FormatUint(uint64(v.Pointer()), 16)
	default:
		return v.Type().String() + " value"
	}
}
