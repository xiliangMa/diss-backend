package system

import (
	"path"

	"github.com/astaxie/beego"
)

const (
	vulnDbRepo = "upload/vuln"
)

// Vulnerability Database Download（漏洞库下载）
type VulnerabilityController struct {
	beego.Controller
}

// @Title DownloadVulnDB
// @Description download vulnerability database file
// @Param token header string true "authorized token"
// @Param If-Modified-Since header string false "vulnerability database last modified timestamp"
// @Param filename path string scanner.db.gz true "vulnerability database filename"
// @Success 200 {object} models.Result
// @router /downloads/vuln/:filename [get]
func (this *VulnerabilityController) DownloadVulnDB() {
	var (
		filename string
	)

	filename = this.GetString(":filename", "scanner.db.gz")
	this.Ctx.Output.Download(path.Join(vulnDbRepo, filename))
}
