package controllers

import (
	"io"
	"net/url"
	"os"

	"github.com/astaxie/beego"
)

type AttachController struct {
	beego.Controller
}

func (ac *AttachController) Get() {
	filePath, err := url.QueryUnescape(ac.Ctx.Request.RequestURI[1:])
	if err != nil {
		ac.Ctx.WriteString(err.Error())
		return
	}

	f, err := os.Open(filePath)
	if err != nil {
		ac.Ctx.WriteString(err.Error())
		return
	}
	defer f.Close()

	_, err = io.Copy(ac.Ctx.ResponseWriter, f)
	if err != nil {
		ac.Ctx.WriteString(err.Error())
		return
	}
}
