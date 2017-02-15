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
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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
	err := http.ListenAndServe(":19999", nil)
	if err != nil {
		log.Fatal("ListenAndServe : ", err)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数,默认不会解析
	contentType := r.Header.Get("contentType")
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		if contentType == "file" {
			w.Write([]byte("Hello,login, here is the POST method and the content type is FILE."))
		} else if contentType == "json" {
			w.Write([]byte("Hello,login, here is the POST method and the content type is JSON."))
		} else if contentType == "form" {
			w.Write([]byte("Hello,login, here is the POST method and the content type is FORM."))
		}
	} else if r.Method == "PUT" { //代理发送的请求居然是put方法
		if contentType == "file" {
			w.Write([]byte("Hello,login, here is the PUT method and the content type is FILE."))
		} else if contentType == "json" {
			w.Write([]byte("Hello,login, here is the PUT method and the content type is JSON."))
		} else if contentType == "form" {
			w.Write([]byte("Hello,login, here is the PUT method and the content type is FORM."))
		}
	} else if r.Method == "DELETE" {
		if contentType == "file" {
			w.Write([]byte("Hello,login, here is the DELETE method and the content type is FILE."))
		} else if contentType == "json" {
			w.Write([]byte("Hello,login, here is the DELETE method and the content type is JSON."))
		} else if contentType == "form" {
			w.Write([]byte("Hello,login, here is the DELETE method and the content type is FORM."))
		}
	}
}

// 处理/upload 逻辑
func upload(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	r.ParseMultipartForm(32 << 20)
	contentType := r.Header.Get("contentType")
	if r.Method == "POST" {
		if contentType == "file" {
			w.Write([]byte("Hello,upload, here is the POST method and the content type is FILE."))
			r.ParseMultipartForm(32 << 20)
			file, handler, err := r.FormFile("uploadfile")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()
			f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666) // 此处假设当前目录下已存在test目录
			if err != nil {
				fmt.Println(err)
				return
			}
			defer f.Close()
			io.Copy(f, file)
		} else if contentType == "json" {
			w.Write([]byte("Hello,upload, here is the POST method and the content type is JSON."))
		} else if contentType == "form" {
			w.Write([]byte("Hello,upload, here is the POST method and the content type is POST."))
		}
	} else if r.Method == "PUT" { //代理发送的请求居然是put方法
		if contentType == "file" {
			w.Write([]byte("Hello,upload, here is the PUT method and the content type is FILE."))
		} else if contentType == "json" {
			w.Write([]byte("Hello,upload, here is the PUT method and the content type is JSON."))
		} else if contentType == "form" {
			w.Write([]byte("Hello,upload, here is the PUT method and the content type is POST."))
		}
	} else if r.Method == "DELETE" {
		if contentType == "file" {
			w.Write([]byte("Hello,upload, here is the DELETE method and the content type is FILE."))
		} else if contentType == "json" {
			w.Write([]byte("Hello,upload, here is the DELETE method and the content type is JSON."))
		} else if contentType == "form" {
			w.Write([]byte("Hello,upload, here is the DELETE method and the content type is POST."))
		}
	}
}
