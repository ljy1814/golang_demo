package controllers

import (
	"blog/models"

	"github.com/astaxie/beego"
)

type ReplyController struct {
	beego.Controller
}

func (rc *ReplyController) Add() {
	tid := rc.Input().Get("tid")
	err := models.AddReply(tid, rc.Input().Get("nickname"), rc.Input().Get("content"))
	if err != nil {
		beego.Error(err)
	}

	rc.Redirect("/topic/view/"+tid, 302)
}

func (rc *ReplyController) Delete() {
	if !checkAccount(rc.Ctx) {
		return
	}
	tid := rc.Input().Get("tid")
	err := models.DeleteReply(rc.Input().Get("rid"))
	if err != nil {
		beego.Error(err)
	}

	rc.Redirect("/topic/view/"+tid, 302)
}
