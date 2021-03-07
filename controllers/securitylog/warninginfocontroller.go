package securitypolicy

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/nats"
	sl "github.com/xiliangMa/diss-backend/service/securitylog"
)

// Warning Info api list
type WarningInfoController struct {
	beego.Controller
}

// @Title GetWarningInfo
// @Description Get Warning Info List (暂不支持租户查询)
// @Param token header string true "authToken"
// @Param body body models.WarningInfo false "告警信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /warninginfo [post]
func (this *WarningInfoController) GetWarningInfoList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	warningInfo := new(models.WarningInfo)
	json.Unmarshal(this.Ctx.Input.RequestBody, &warningInfo)

	this.Data["json"] = warningInfo.List(from, limit)
	this.ServeJSON(false)
}

// @Title UpdateWarningInfo
// @Description Update WarningInfo
// @Param token header string true "authToken"
// @Param id path string "" true "id"
// @Param body body models.WarningInfo false "WarningInfo"
// @Success 200 {object} models.Result
// @router /warninginfo/:id [put]
func (this *WarningInfoController) UpdateWarningInfo() {
	id := this.GetString(":id")
	warningInfo := new(models.WarningInfo)
	json.Unmarshal(this.Ctx.Input.RequestBody, &warningInfo)

	warningInfo.Id = id
	result := warningInfo.Update()

	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title AddClientSub_Image_Safe
// @Description Add Nats Subscribe topic for Image_Safe
// @Param token header string true "authToken"
// @Param libname path string "" true "Registry Name"
// @Success 200 {object} models.Result
// @router /warninginfo/add_sub_imagelib/:libname [post]
func (this *WarningInfoController) AddClientSub_Image_Safe() {
	libname := this.GetString(":libname")
	result := nats.AddClientSub_Image_Safe(libname)

	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title WarninginfoDisposal
// @Description Add Alarm Processing
// @Param token header string true "authToken"
// @Param body body models.WarningInfo true "WarningInfo"
// @Success 200 {object} models.Result
// @router /warninginfo/disposal [post]
func (this *WarningInfoController) DisposalMode() {
	warningInfo := new(models.WarningInfo)
	json.Unmarshal(this.Ctx.Input.RequestBody, &warningInfo)
	wis := new(sl.WarningInfoService)
	result := wis.DisposalMode(warningInfo)

	this.Data["json"] = result
	this.ServeJSON(false)
}
