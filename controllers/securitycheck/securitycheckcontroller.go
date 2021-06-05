package securitycheck

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
	securitycheck "github.com/xiliangMa/diss-backend/service/securitycheck"
	"time"
)

// 安全检测接口列表
type SecurityCheckController struct {
	beego.Controller
}

// @Title SecurityCheck
// @Description Security check
// @Param token header string true "authToken"
// @Param account query string "admin" false "租户"
// @Param body body models.SecurityCheckList true "检查列表"
// @Success 200 {object} models.Result
// @router / [post]
func (this *SecurityCheckController) SecurityCheck() {
	checkList := new(models.SecurityCheckList)
	account := this.GetString("account")
	if account == "" {
		account = models.Account_Admin
	}
	json.Unmarshal(this.Ctx.Input.RequestBody, &checkList)
	batch := time.Now().UnixNano() / 1e3
	securityCheckService := securitycheck.SecurityCheckService{SecurityCheckList: checkList, Batch: batch}
	this.Data["json"] = securityCheckService.DeliverTask()
	this.ServeJSON(false)
}

// @Title SecurityCheck2
// @Description Security Check2
// @Param token header string true "authToken"
// @Param isSystem query bool false false "系统检查"
// @Param body body models.SecurityCheckParams true "安全检查参数"
// @Success 200 {object} models.Result
// @router /v2 [post]
func (this *SecurityCheckController) SecurityCheck2() {
	params := new(models.SecurityCheckParams)
	isSystem, _ := this.GetBool("isSystem")
	json.Unmarshal(this.Ctx.Input.RequestBody, &params)
	batch := time.Now().UnixNano()
	securityScanService := securitycheck.SecurityScanService{SecurityCheckParams: params, Batch: batch, IsSystem: isSystem}
	this.Data["json"] = securityScanService.DeliverTask()
	this.ServeJSON(false)
}
