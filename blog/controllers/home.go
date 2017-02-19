package controllers

import (
	"blog/models"
	"fmt"

	"github.com/astaxie/beego"
)

type HomeController struct {
	beego.Controller
}

func (hc *HomeController) Get() {
	fmt.Println("xxxx")
	hc.Data["IsHome"] = true
	hc.TplName = "home.html"
	hc.Data["IsLogin"] = checkAccount(hc.Ctx)

	topics, err := models.GetAllTopics(hc.Input().Get("cate"), hc.Input().Get("label"), true)
	if err != nil {
		beego.Error(err)
	}
	hc.Data["Topics"] = topics

	categories, err := models.GetAllCategories()
	if err != nil {
		beego.Error(err)
	}
	hc.Data["Categories"] = categories
}
