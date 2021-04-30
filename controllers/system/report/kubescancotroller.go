package report

import (
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/service/kubevuln"
)

// Kube vuln log api
type KubeVulnLogController struct {
	beego.Controller
}

// @Title KubeVulnScanLog
// @Description KubeVulnScanLog
// @Param body body string false "集群漏洞扫描结果"
// @Success 200 {object} models.Result
// @router /kubevulnscan [post]
func (this *KubeVulnLogController) ReceiveKubeVulnScanLog() {
	rawLog := this.Ctx.Input.RequestBody
	kubeVlunService := kubevuln.KubeVlunService{RawLog: rawLog}
	kubeVlunService.ReceiveKubeScanLog()
	this.ServeJSON(false)
}
