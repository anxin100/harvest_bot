package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	"net/http"
)

type ChatController struct {
	MainController
}

// @router /receive [post]
func (server *ChatController) Receive() {
	body := server.Ctx.Input.RequestBody
	logs.Debug("Receive msg body = ", string(body))

	param := server.Ctx.Input.Params()
	logs.Debug("Receive msg body = ", param)

	data := server.Ctx.Input.Data()
	logs.Debug("Receive msg body = ", data)

	server.respond(http.StatusOK, "我收到了你的消息")
}
