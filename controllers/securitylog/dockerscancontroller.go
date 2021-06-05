package securitypolicy

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// docker scan log api list
type DockerScanController struct {
	beego.Controller
}

// @Title GetDockerScan
// @Description Get DockerScan Log
// @Param token header string true "authToken"
// @Param body body models.DockerVulnerabilities false "docker漏洞扫描结果"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /dockerscan [post]
func (this *DockerScanController) GetDockerScan() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	dockerScan := new(models.DockerVulnerabilities)
	json.Unmarshal(this.Ctx.Input.RequestBody, &dockerScan)

	this.Data["json"] = dockerScan.List(from, limit)
	this.ServeJSON(false)
}
