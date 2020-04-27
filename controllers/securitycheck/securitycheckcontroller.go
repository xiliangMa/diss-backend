package securitycheck

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models/bean"
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
// @Param body body bean.SecurityCheckList true "检查列表"
// @Param nats query bool "false" false "是否下发给nats"
// @Success 200 {object} models.Result
// @router / [post]
func (this *SecurityCheckController) SecurityCheck() {
	checkList := new(bean.SecurityCheckList)
	isNats, _ := this.GetBool("nats")
	json.Unmarshal(this.Ctx.Input.RequestBody, &checkList)
	bath := time.Now().Unix()
	securityCheckService := securitycheck.SecurityCheckService{checkList, nil, bath, nil}
	this.Data["json"] = securityCheckService.DeliverTask(isNats)
	this.ServeJSON(false)
}
