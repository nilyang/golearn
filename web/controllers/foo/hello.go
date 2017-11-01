package foo

import "github.com/astaxie/beego"

type HelloController struct{
	beego.Controller
}

func (f *HelloController) Hello()  {

	var result = map[string]interface{} {
		"foo":1,
		"bar":2,
		"appName":beego.BConfig.AppName,
		"foox": beego.AppConfig.String("foo"),
	}
	f.Data["json"] = result
	f.ServeJSON()
}

