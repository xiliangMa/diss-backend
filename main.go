package main

import (
	"github.com/astaxie/beego"
	_ "github.com/xiliangMa/diss-backend/routers"
	_ "github.com/xiliangMa/diss-backend/sysinit"
)

func main() {
	beego.BConfig.WebConfig.DirectoryIndex = true
	beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	beego.Run()
}
