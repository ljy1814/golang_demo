package main

import (
	"bee01/controllers"
	"bee01/models"
	"io"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type HomeController struct {
	beego.Controller
}

func (h *HomeController) Get() {
	h.Ctx.WriteString("Hello, World!")
}

func main() {
	//	beego.Router("/", &HomeController{})
	//	beego.Run(":1234")

	/*
		http.HandleFunc("/hello", sayHello)
		http.HandleFunc("/bye", sayBye)

		err := http.ListenAndServe(":1234", nil)
		if err != nil {
			log.Fatal(err)
		}
	*/

	orm.Debug = true
	//自动建表,自动在sqlite3数据库建立表结构
	orm.RunSyncdb("default", false, true)
	beego.Router("/", &controllers.MainController{})
	beego.Router("/category", &controllers.CategoryController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/home", &controllers.HomeController{})
	beego.Router("/topic", &controllers.TopicController{})
	beego.AutoRouter(&controllers.TopicController)
	beego.Run(":2345")
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello golang, welcome join us.")
}

func sayBye(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Bye bye, ")
}

func init() {
	//注册数据库
	models.RegisterDB()
}
