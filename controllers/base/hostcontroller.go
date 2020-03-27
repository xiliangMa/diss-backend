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
