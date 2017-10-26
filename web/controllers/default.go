package controllers

import (
	"github.com/astaxie/beego"
	"golearn/web/models"
)

type MainController struct {
	beego.Controller
}

var hub *models.Hub

func init() {
	hub = models.NewHub()
	go hub.Run()
}


func (c *MainController) Get() {
	c.TplName = "index.html"
}


func (c *MainController) WsGet() {
	c.EnableRender = false
	models.ServeWebsocket(hub, c.Ctx.ResponseWriter, c.Ctx.Request)
}

func (c *MainController) SendMsg() {
	c.EnableRender = false
	msg := c.GetString("msg")
	if msg == "" {
		msg = "Hello world!"
	}
	models.SendMessage(hub,msg)
}
