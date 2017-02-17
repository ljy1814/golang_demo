package controllers

import "github.com/astaxie/beego"

type TopicController struct {
	beego.Controller
}

func (tc *TopicController) Get() {
	tc.Data["IsTopic"] = true
	tc.TplName = "topic.html"
	tc.Data["IsLogin"] = checkAccount(false)

	topics, err := models.GetAllTopics(false)
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

	var err error
	if len(tid) == 0 {
		err = models.AddTopics(title, content)
	} else {
		err = models.ModiifyTopic(tid, title, content)
	}

	if err != nil {
		begoo.Error(err)
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
}

func (tc *TopicController) Add() {
	if !checkAccount(tc.Ctx) {
		tc.Redirect("/login", 302)
		return
	}

	tc.TplName = "topic_add.html"
}

func (tc *TopicController) Delete() {
	if !checkAccount(tc.Ctx) {
		tc.Redirect("/login", 302)
		return
	}

	err := models.DeleteTopic(tc.Input().get("tid"))
	if err != nil {
		beego.Error(err)
	}

	tc.Redirect("/topic", 302)
}

func (tc *TopicController) View() {
	tc.TplName = "topic_view.html"

	topic, err := models.GetTopic(tc.Ctx.Input.Params["0"])
	if err != nil {
		beego.Error(err)
		tc.Redirect("/", 302)
		return
	}

	tc.Data["Topic"] = topic
}
