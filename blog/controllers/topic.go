package controllers

import (
	"blog/models"
	"path"
	"strings"

	"github.com/astaxie/beego"
)

type TopicController struct {
	beego.Controller
}

func (tc *TopicController) Get() {
	tc.Data["IsTopic"] = true
	tc.TplName = "topic.html"
	tc.Data["IsLogin"] = checkAccount(tc.Ctx)

	topics, err := models.GetAllTopics("", "", false)
	if err != nil {
		beego.Error(err)
	}
	tc.Data["Topics"] = topics
}

func (tc *TopicController) Post() {
	if !checkAccount(tc.Ctx) {
		tc.Redirect("/login", 302)
		return
	}

	tid := tc.Input().Get("tid")
	title := tc.Input().Get("title")
	content := tc.Input().Get("content")
	category := tc.Input().Get("category")
	label := tc.Input().Get("label")

	_, fh, err := tc.GetFile("attachment")
	if err != nil {
		beego.Error(err)
	}
	var attachment string
	if fh != nil {
		attachment = fh.Filename
		beego.Info(attachment)
		//保持附件
		err = tc.SaveToFile("attachment", path.Join("attachment", attachment))
		if err != nil {
			beego.Error(err)
		}
	}

	if len(tid) == 0 {
		err = models.AddTopic(title, category, label, content, attachment)
	} else {
		err = models.ModifyTopic(tid, title, category, label, content, attachment)
	}

	if err != nil {
		beego.Error(err)
	}

	tc.Redirect("/topic", 302)

}

func (tc *TopicController) Modify() {
	tc.TplName = "topic_modify.html"
	tid := tc.Input().Get("tid")
	topic, err := models.GetTopic(tid)
	if err != nil {
		beego.Error(err)
		tc.Redirect("/", 302)
		return
	}
	tc.Data["Topic"] = topic
	tc.Data["Tid"] = tid
	tc.Data["IsLogin"] = true
}

func (tc *TopicController) Add() {
	if !checkAccount(tc.Ctx) {
		tc.Redirect("/login", 302)
		return
	}

	tc.TplName = "topic_add.html"
	tc.Data["IsLogin"] = true
}

func (tc *TopicController) Delete() {
	if !checkAccount(tc.Ctx) {
		tc.Redirect("/login", 302)
		return
	}

	err := models.DeleteTopic(tc.Input().Get("tid"))
	if err != nil {
		beego.Error(err)
	}

	tc.Redirect("/topic", 302)
}

func (tc *TopicController) View() {
	tc.TplName = "topic_view.html"

	//	paths := strings.Split(tc.Ctx.Request.URL.Path, "/")
	//	tid := paths[len(paths)-1]
	reqUrl := tc.Ctx.Request.RequestURI
	i := strings.LastIndex(reqUrl, "/")
	tid := reqUrl[i+1:]
	topic, err := models.GetTopic(tid)
	if err != nil {
		beego.Error(err)
		tc.Redirect("/", 302)
		return
	}

	tc.Data["Topic"] = topic
	tc.Data["Labels"] = strings.Split(topic.Labels, " ")
	replies, err := models.GetAllReplies(tid)
	if err != nil {
		beego.Error(err)
		return
	}

	tc.Data["Replies"] = replies
	tc.Data["IsLogin"] = checkAccount(tc.Ctx)
}
