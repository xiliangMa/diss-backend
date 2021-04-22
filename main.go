package main

import (
	"github.com/astaxie/beego"
	_ "github.com/xiliangMa/diss-backend/routers"
	_ "github.com/xiliangMa/diss-backend/sysinit"
	"github.com/xiliangMa/diss-backend/utils"
)

func main() {
	// logo
	logoPath := utils.GetLogoPath() + utils.GetLogoName()
	logoUrl := utils.GetLogoUrl()
	beego.SetStaticPath(logoUrl, logoPath)
	// vuln
	beego.SetStaticPath(utils.GetVulnDbUrl(), utils.GetVulnDbPath())
	// probe driver
	beego.SetStaticPath(utils.GetProbeDriverUrl(), utils.GetProbeDriverPath())
	beego.BConfig.WebConfig.DirectoryIndex = true
	beego.BConfig.WebConfig.StaticDir["/"] = "swagger"
	beego.Run()
}
