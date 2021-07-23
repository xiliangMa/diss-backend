package k8s

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// Deployment接口
type DeploymentController struct {
	beego.Controller
}

// @Title GetDeployment
// @Description Get Deployment List
// @Param token header string true "authToken"
// @Param body body models.Deployment false "Deployment"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *DeploymentController) GetDeploymentList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	Deployment := new(models.Deployment)
	json.Unmarshal(this.Ctx.Input.RequestBody, &Deployment)
	this.Data["json"] = Deployment.List(from, limit)
	this.ServeJSON(false)
}
