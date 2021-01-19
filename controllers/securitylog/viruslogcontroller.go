package securitypolicy

import (
	"encoding/json"
	"github.com/beego/beego/v2/server/web"
	"github.com/xiliangMa/diss-backend/models"
)

// Virus Log api list
type VirusLogController struct {
	web.Controller
}

// @Title GetImageVirusLog
// @Description Get ImageVirus Log List
// @Param token header string true "authToken"
// @Param body body models.ImageVirus false "镜像病毒信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /virus/image [post]
func (this *VirusLogController) GetVirusLogList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	imageVirus := new(models.ImageVirus)
	json.Unmarshal(this.Ctx.Input.RequestBody, &imageVirus)
	this.Data["json"] = imageVirus.List(from, limit)
	this.ServeJSON(false)
}

// @Title GetHostOrContainerVirusLog
// @Description Get HostOrContainer Virus Log List (1. 根据 TargeType = host 和 HostId = All 判断是否是查询所有主机日志 如果不是则匹配其它所传入的条件 2. 根据 TargeType = container 和 ContainerId = All 判断是否是查询所有容器日志 如果不是则匹配其它所传入的条件)
// @Param token header string true "authToken"
// @Param body body models.DockerVirus false "主机或者容器病毒信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /virus/hostorcontainer [post]
func (this *VirusLogController) GetHostOrContainerVirusLogList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	dockerVirus := new(models.DockerVirus)
	json.Unmarshal(this.Ctx.Input.RequestBody, &dockerVirus)
	this.Data["json"] = dockerVirus.List(from, limit)
	this.ServeJSON(false)
}
