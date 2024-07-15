package main

import (
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
	_ "harvest_bot/routers"
)

func init() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Async", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "X-Xsrf-Token"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "X-Xsrf-Token", "Authorization", "Async"},
		AllowCredentials: true,
	}))
	timeout := beego.BConfig.Listen.ServerTimeOut
	logs.Debug("timeout = ", timeout)
}

func main() {
	//beego.BConfig.WebConfig.AutoRender = false
	//beego.BConfig.WebConfig.EnableDocs = true // 启用路由调试
	beego.Run()
}
