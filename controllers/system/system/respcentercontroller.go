package system

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
	css "github.com/xiliangMa/diss-backend/service/system/system"
)

// Container Response Center api list
type RespCenterController struct {
	beego.Controller
}

// @Title GetRespCenterList
// @Description Get Resp Center List (暂不支持租户查询)
// @Param token header string true "authToken"
// @Param body body models.RespCenter false "响应中心"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /respcenter [post]
func (this *RespCenterController) GetRespCenterList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	crc := new(models.RespCenter)
	json.Unmarshal(this.Ctx.Input.RequestBody, &crc)

	this.Data["json"] = crc.List(from, limit)
	this.ServeJSON(false)
}

// @Title GetRespCenter
// @Description Get Resp Center (暂不支持租户查询)
// @Param token header string true "authToken"
// @Param id path string "" true "id"
// @Success 200 {object} models.Result
// @router /respcenter/:id [post]
func (this *RespCenterController) GetRespCenter() {
	id := this.GetString(":id")
	rc := new(models.RespCenter)
	json.Unmarshal(this.Ctx.Input.RequestBody, &rc)

	rc.Id = id
	this.Data["json"] = rc.GetRespCenter()
	this.ServeJSON(false)
}

// @Title ContainerOperation
// @Description Resp Center operation
// @Param token header string true "authToken"
// @Param body body models.RespCenter true "RespCenter"
// @Success 200 {object} models.Result
// @router /respcenter/operation [post]
func (this *RespCenterController) ContainerOperation() {
	resp := new(models.RespCenter)
	json.Unmarshal(this.Ctx.Input.RequestBody, &resp)
	respCenterService := new(css.RespCenterService)
	result := respCenterService.ContainerOperation(resp)

	this.Data["json"] = result
	this.ServeJSON(false)
}
