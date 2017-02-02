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
	} else {
		r.ParseForm() //解析参数,默认不会解析
		fmt.Println("username : ", r.Form["username"])
		fmt.Println("password : ", r.Form["password"])
		fmt.Println(r.Form)
	}
}
