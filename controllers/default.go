package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

func (server *MainController) respond(code int, message string, data ...interface{}) {
	status := 200
	if code == 401 {
		status = code
	}
	server.Ctx.Output.SetStatus(status)
	var d interface{}
	if len(data) > 0 {
		d = data[0]
	}
	server.Data["json"] = struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}{
		Code:    code,
		Message: message,
		Data:    d,
	}
	server.ServeJSON()
}
