package main

import (
	"blog/controllers"
	"blog/models"
	_ "blog/routers"
	"os"
	"reflect"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func main() {
	orm.Debug = true
	//自动建表,自动在sqlite3数据库建立表结构
	orm.RunSyncdb("default", false, true)
	beego.Router("/", &controllers.HomeController{})
	beego.Router("/category", &controllers.CategoryController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/home", &controllers.HomeController{})
	beego.Router("/topic", &controllers.TopicController{})
	beego.AutoRouter(&controllers.TopicController{})
	beego.Router("/reply", &controllers.ReplyController{})
	beego.Router("/reply/add", &controllers.ReplyController{}, "post:Add")
	beego.Router("/reply/delete", &controllers.ReplyController{}, "get:Delete")

	//添加附件
	os.Mkdir("attachment", os.ModePerm)
	beego.Router("/attachment/:all", &controllers.AttachController{})

	beego.AddFuncMap("typ", getType)
	beego.Run(":2345")
}

func getType(ss []string) string {
	return reflect.TypeOf(ss).String()
}

func init() {
	//注册数据库
	models.RegisterDB()
}
