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
	fmt.Println("method : ", r.Method)  //获取请求方法
	fmt.Println("path", r.URL.Path)     //
	fmt.Println("host", r.URL.Host)     //
	fmt.Println("schema", r.URL.Scheme) //
	fmt.Println("RawQuery", r.URL.RawQuery) //
	r.ParseForm() //解析参数,默认不会解析
	fmt.Println(r.Form)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else if r.Method == "POST" {
	fmt.Println(r.Form["type"][0])
        if r.Form["type"][0] == "file" {
            w.Write([]byte("Hello,login, here is the POST method and the content type is FILE."))
        } else if r.Form["type"][0] == "json" {
            w.Write([]byte("Hello,login, here is the POST method and the content type is JSON."))
        } else if r.Form["type"][0] == "form" {
            w.Write([]byte("Hello,login, here is the POST method and the content type is FORM."))
        }
    } else if r.Method == "PUT" { //代理发送的请求居然是put方法
        if r.Form["type"][0] == "file" {
            w.Write([]byte("Hello,login, here is the PUT method and the content type is FILE."))
        } else if r.Form["type"][0] == "json" {
            w.Write([]byte("Hello,login, here is the PUT method and the content type is JSON."))
        } else if r.Form["type"][0] == "form" {
            w.Write([]byte("Hello,login, here is the PUT method and the content type is FORM."))
        }
    } else if r.Method == "DELETE" {
        if r.Form["type"][0] == "file" {
            w.Write([]byte("Hello,login, here is the DELETE method and the content type is FILE."))
        } else if r.Form["type"][0] == "json" {
            w.Write([]byte("Hello,login, here is the DELETE method and the content type is JSON."))
        } else if r.Form["type"][0] == "form" {
            w.Write([]byte("Hello,login, here is the DELETE method and the content type is FORM."))
        }
    }
}

// 处理/upload 逻辑
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)        //获取请求的方法
	fmt.Println("path", r.URL.Path)         //
	fmt.Println("schema", r.URL.Scheme)     //
	fmt.Println("Path", r.URL.Path)         //
	fmt.Println("RawQuery", r.URL.RawQuery) //
	fmt.Println(r.Form["url_long"])
	r.ParseForm()
	if r.Method == "POST" {
		if r.Form["type"][0] == "file" {
			w.Write([]byte("Hello,upload, here is the POST method and the content type is FILE."))
		} else if r.Form["type"][0] == "json" {
			w.Write([]byte("Hello,upload, here is the POST method and the content type is JSON."))
		} else if r.Form["type"][0] == "form" {
			w.Write([]byte("Hello,upload, here is the POST method and the content type is POST."))
		}
	} else if r.Method == "PUT" { //代理发送的请求居然是put方法
		if r.Form["type"][0] == "file" {
			w.Write([]byte("Hello,upload, here is the PUT method and the content type is FILE."))
		} else if r.Form["type"][0] == "json" {
			w.Write([]byte("Hello,upload, here is the PUT method and the content type is JSON."))
		} else if r.Form["type"][0] == "form" {
			w.Write([]byte("Hello,upload, here is the PUT method and the content type is POST."))
		}
	} else if r.Method == "DELETE" {
		if r.Form["type"][0] == "file" {
			w.Write([]byte("Hello,upload, here is the DELETE method and the content type is FILE."))
		} else if r.Form["type"][0] == "json" {
			w.Write([]byte("Hello,upload, here is the DELETE method and the content type is JSON."))
		} else if r.Form["type"][0] == "form" {
			w.Write([]byte("Hello,upload, here is the DELETE method and the content type is POST."))
		}
    }
        
}
