package main

import (
	"github.com/astaxie/beego"
	_ "github.com/xiliangMa/diss-backend/routers"
	_ "github.com/xiliangMa/diss-backend/sysinit"
)

func main() {
	logoParh := beego.AppConfig.String("system::LogoPath") + beego.AppConfig.String("system::NewLogoName")
	logoUrl := beego.AppConfig.String("system::LogoUrl")
	beego.SetStaticPath(logoUrl, logoParh)
	beego.BConfig.WebConfig.DirectoryIndex = true
	beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	beego.Run()
}
