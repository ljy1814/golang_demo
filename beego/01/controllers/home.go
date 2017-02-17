package controllers

import "github.com/astaxie/beego"

type HomeController struct {
	beego.Controller
}

func (hc *HomeController) Get() {
	hc.Data["IsHome"] = true
	hc.TplName = "home.html"
	hc.Data["IsLogin"] = checkAccount(hc.Ctx)
}
