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
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
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
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		//		r.ParseForm() //解析参数,默认不会解析
		nUrl := handleUrl(r)
		r.ParseForm()
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
		nUrl := handleUrl(r)
		r.ParseMultipartForm(32 << 20)
		_, mff, err := r.FormFile("uploadfile")
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
}

func handleUrl(r *http.Request) string {
	url := r.URL
	nUrl := newUrl(url)
	if strings.Index(nUrl, "?") == len(nUrl)-1 {
		nUrl += "channel=" + fmt.Sprintf("%x", md5.New().Sum(nil))
	} else {
		nUrl += "&channel=" + fmt.Sprintf("%x", md5.New().Sum(nil))
	}
	nUrl += "&sign=" + genSign(url)
	return nUrl
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
