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
// @Description Security heck
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
	batch := time.Now().Unix()
	securityCheckService := securitycheck.SecurityCheckService{SecurityCheckList: checkList, Batch: batch, Account: account}
	this.Data["json"] = securityCheckService.DeliverTask()
	this.ServeJSON(false)
}
