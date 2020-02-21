package asset

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

// @Title GetNameSpaceList
// @Description Get NameSpace List
// @Param token header string true "auth token"
// @Param clusterId path string "" true "clusterId"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /k8s/clusters/:clusterId/namespaces [post]
func (this *K8SController) GetNameSpaces() {
	clusterId := this.GetString(":clusterId")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	ns := new(k8s.NameSpace)
	ns.ClusterId = clusterId
	this.Data["json"] = ns.List(from, limit)
	this.ServeJSON(false)

}

// @Title GetPodList
// @Description Get Pod List
// @Param token header string true "auth token"
// @Param nsName path string "" true "namespaceName"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /k8s/namespaces/:nsName/pods [post]
func (this *K8SController) GetPods() {
	nsName := this.GetString(":nsName")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	pod := new(k8s.Pod)
	pod.NameSpaceName = nsName
	this.Data["json"] = pod.List(from, limit)
	this.ServeJSON(false)

}
