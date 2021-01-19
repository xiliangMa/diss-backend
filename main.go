package main

import (
	"github.com/beego/beego/v2/server/web"
	_ "github.com/xiliangMa/diss-backend/routers"
	_ "github.com/xiliangMa/diss-backend/sysinit"
)

func main() {
	logoPath, _ := web.AppConfig.String("system::LogoPath")
	logoName, _ := web.AppConfig.String("system::NewLogoName")
	path := logoPath + logoName
	logoUrl, _ := web.AppConfig.String("system::LogoUrl")
	web.SetStaticPath(logoUrl, path)
	web.BConfig.WebConfig.DirectoryIndex = true
	web.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	web.Run()
}
