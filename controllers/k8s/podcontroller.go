package k8s

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models/k8s"
)

// Pod接口
type PodController struct {
	beego.Controller
}

// @Title GetPod
// @Description Get Pod List
// @Param token header string true "authToken"
// @Param body body k8s.Pod false "Pod"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *PodController) GetPodsList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	pod := new(k8s.Pod)
	json.Unmarshal(this.Ctx.Input.RequestBody, &pod)
	this.Data["json"] = pod.List(from, limit)
	this.ServeJSON(false)
}
