package routers

import (
	"github.com/astaxie/beego"
	"golearn/web/controllers"
	"golearn/web/controllers/foo"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("/ws", &controllers.MainController{},"get:WsGet")
	beego.Router("/send", &controllers.MainController{},"get:SendMsg")
	beego.Router("/hello",&foo.HelloController{},"get:Hello")
}
