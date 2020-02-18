package k8s

import (
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models/k8s"
)

// Asset K8S object api list
type K8SController struct {
	beego.Controller
}

// @Title GetClusterList
// @Description Get Cluster List
// @Param token header string true "auth token"
// @Param name query string "" false "name"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /k8s/clusters [post]
func (this *K8SController) GetClusterList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	cluster := new(k8s.Cluster)
	cluster.Name = this.GetString("name")
	this.Data["json"] = cluster.List(from, limit)
	this.ServeJSON(false)

}
