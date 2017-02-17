package controllers

import (
	"bee01/models"

	"github.com/astaxie/beego"
)

type CategoryController struct {
	beego.Controller
}

func (cc *CategoryController) Get() {
	op := cc.Input().Get("op")
	switch op {
	case "add":
		name := cc.Input().Get("name")
		if len(name) == 0 {
			break
		}

		err := models.AddCategory(name)
		if err != nil {
			beego.Error(err)
		}

		cc.Redirect("/category", 302)
		return
	case "del":
		id := cc.Input().Get("id")
		if len(id) == 0 {
			break
		}

		err := models.DeleteCategory(id)
		if err != nil {
			beego.Error(err)
		}

		cc.Redirect("/category", 302)
		return
	}

	cc.Data["IsCategory"] = true
	cc.TplName = "category.html"
	cc.Data["IsLogin"] = checkAccount(cc.Ctx)

	var err error
	cc.Data["Categories"], err = models.GetAllCategories()
	if err != nil {
		beego.Error(err)
	}
}
