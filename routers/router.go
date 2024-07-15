package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"harvest_bot/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.AutoPrefix("app/api/v1", &controllers.ChatController{})

	//ns := beego.NewNamespace("/app",
	//	beego.NSNamespace("/api",
	//		beego.NSNamespace("/v1",
	//			beego.NSNamespace("/chat",
	//				beego.NSInclude(
	//					&controllers.ChatController{},
	//				),
	//			),
	//		),
	//	),
	//)
	//beego.AddNamespace(ns)
}
