package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

    beego.GlobalControllerRouter["harvest_bot/controllers:ChatController"] = append(beego.GlobalControllerRouter["harvest_bot/controllers:ChatController"],
        beego.ControllerComments{
            Method: "Receive",
            Router: "/receive",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
