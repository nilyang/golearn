package routers

import (
	"golearn/web/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("/ws", &controllers.MainController{},"get:WsGet")
	beego.Router("/send", &controllers.MainController{},"get:SendMsg")
}
