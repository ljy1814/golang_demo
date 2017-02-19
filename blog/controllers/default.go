package controllers

import "github.com/astaxie/beego"

type MainController struct {
	beego.Controller
}

func (mc *MainController) Get() {
	//若是直接写入会跳过模板渲染
	//	mc.Ctx.WriteString("AppName : " + beego.AppConfig.String("appname") + "\nRunMode : " + beego.AppConfig.String("runmode"))

	//	mc.Ctx.WriteString("\n\nAppName : " + beego.Appname + "\nRunMode : " + beego.RunMode)
	mc.Data["Username"] = "傻屌"
	mc.Data["Email"] = "yajin@shaodiao.shop"
	mc.TplName = "index.tpl"

	mc.Data["TrueCond"] = true
	mc.Data["FalseCond"] = false

	type u struct {
		Name string
		Age  int
		Sex  string
	}

	user := u{
		Name: "傻屌",
		Age:  20,
		Sex:  "Male",
	}

	mc.Data["User"] = user

	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	mc.Data["Nume"] = nums
	users := make([]u, 0, 1)
	users = append(users, user)
	mc.Data["Users"] = users

	mc.Data["TplVar"] = "hello, 小傻"

	mc.Data["Html"] = "<div>hello beego</div>"

	mc.Data["Pipe"] = "<div>The string will bve escaped through pipeline and template function function</div>"

	beego.Trace("Trace test1")
	beego.Info("Info test1")
	//	beego.SetLevel(beego.LevelInfo)

	beego.Trace("Trace test2")
	beego.Info("Info test2")
}
