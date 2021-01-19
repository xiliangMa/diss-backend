package controllers

import "github.com/beego/beego/v2/server/web"

type NestPreparer interface {
	NestPrepare()
}

type baseController struct {
	web.Controller
}

func (this *baseController) Prepare() {
	if app, ok := this.AppController.(NestPreparer); ok {
		app.NestPrepare()
	}
}
