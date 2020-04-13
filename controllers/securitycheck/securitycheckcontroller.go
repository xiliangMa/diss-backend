package securitycheck

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models/bean"
)

// 安全检测接口列表
type SecurityCheckController struct {
	beego.Controller
}

// @Title SecurityCheck
// @Description Security heck
// @Param token header string true "authToken"
// @Param hostId path string "" true "hostId"C
// @Param body body bean.SecurityCheckList true "检查列表"
// @Success 200 {object} models.Result
// @router / [post]
func (this *SecurityCheckController) SecurityCheck() {
	checkList := new(bean.SecurityCheckList)
	json.Unmarshal(this.Ctx.Input.RequestBody, &checkList)
	// to do
	this.Data["json"] = nil
	this.ServeJSON(false)
}
