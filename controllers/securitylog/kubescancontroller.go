package securitypolicy

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// kube scan log api list
type KubeScanController struct {
	beego.Controller
}

// @Title GetKubeScan
// @Description Get kubescan Log
// @Param token header string true "authToken"
// @Param body body models.KubeScan false "集群漏洞扫描结果"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /kubecsan [post]
func (this *KubeScanController) GetKubeScan() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	kubeScan := new(models.KubeScan)
	json.Unmarshal(this.Ctx.Input.RequestBody, &kubeScan)
	this.Data["json"] = kubeScan.List(from, limit)
	this.ServeJSON(false)
}
