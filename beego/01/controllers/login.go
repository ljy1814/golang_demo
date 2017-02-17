package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

type LoginController struct {
	beego.Controller
}

func (lc *LoginController) Get() {
	if lc.Input().Get("exit") == "true" {
		lc.Ctx.SetCookie("uname", "", -1, "/")
		lc.Ctx.SetCookie("pwd", "", -1, "/")
		lc.Redirect("/", 302)
		return
	}
	lc.TplName = "login.html"
}

func (lc *LoginController) Post() {
	uname := lc.Input().Get("uname")
	pwd := lc.Input().Get("pwd")

	autoLogin := lc.Input().Get("autoLogin") == "on"

	if uname == beego.AppConfig.String("adminName") && pwd == beego.AppConfig.String("adminPass") {
		maxAge := 0
		if autoLogin {
			maxAge = 1<<31 - 1
		}

		lc.Ctx.SetCookie("uname", uname, maxAge, "/")
		lc.Ctx.SetCookie("pwd", pwd, maxAge, "/")
	}

	lc.Redirect("/", 302)
	return
}

func checkAccount(ctx *context.Context) bool {
	ck, err := ctx.Request.Cookie("uname")
	if err != nil {
		return false
	}

	uname := ck.Value

	ck, err = ctx.Request.Cookie("pwd")
	if err != nil {
		return false
	}

	pwd := ck.Value
	return uname == beego.AppConfig.String("adminName") && pwd == beego.AppConfig.String("adminPass")
}
