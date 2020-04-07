package base

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// 主机接口列表
type HostController struct {
	beego.Controller
}

// @Title GetHosts
// @Description Get Hosts
// @Param token header string true "authToken"
// @Param body body models.HostConfig false "主机配置信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *HostController) HostList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	hostConfig := new(models.HostConfig)
	json.Unmarshal(this.Ctx.Input.RequestBody, &hostConfig)
	this.Data["json"] = hostConfig.List(from, limit)
	this.ServeJSON(false)
}

// @Title HostPs
// @Description Get HostPs List
// @Param token header string true "authToken"
// @Param hostId path string "" true "hostId"
// @Param body body models.HostPs false "主机进程"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /:hostId/ps [post]
func (this *HostController) GetHostPsList() {
	hostId := this.GetString(":hostId")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	hostPs := new(models.HostPs)
	json.Unmarshal(this.Ctx.Input.RequestBody, &hostPs)
	hostPs.HostId = hostId
	this.Data["json"] = hostPs.List(from, limit)
	this.ServeJSON(false)
}

// @Title UpdateHost
// @Description Update Host
// @Param token header string true "authToken"
// @Param hostId path string "" true "hostId"
// @Param body body models.HostConfig true "主机配置信息"
// @Success 200 {object} models.Result
// @router /:hostId [put]
func (this *HostController) UpdateHost() {
	hostId := this.GetString(":hostId")
	hostConfig := new(models.HostConfig)
	json.Unmarshal(this.Ctx.Input.RequestBody, &hostConfig)
	hostConfig.Id = hostId
	this.Data["json"] = hostConfig.Update()
	this.ServeJSON(false)
}
