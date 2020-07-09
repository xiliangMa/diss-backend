package securitypolicy

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
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
