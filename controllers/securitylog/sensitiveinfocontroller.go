package securitypolicy

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// SensitiveInfo Log api
type SensitiveInfoController struct {
	beego.Controller
}

// @Title GetSensitiveInfo
// @Description Get SensitiveInfo List
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param body body models.SensitiveInfo false "敏感信息列表"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /sensitiveinfo [post]
func (this *VirusLogController) GetSensitiveInfoList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	sensiInfo := new(models.SensitiveInfo)
	json.Unmarshal(this.Ctx.Input.RequestBody, &sensiInfo)

	this.Data["json"] = sensiInfo.List(from, limit)
	this.ServeJSON(false)
}
