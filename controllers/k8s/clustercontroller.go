package k8s

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// 集群接口
type ClusterController struct {
	beego.Controller
}

// @Title GetClusters
// @Description Get Cluster List(不支持租户查询)
// @Param token header string true "authToken"
// @Param body body models.Cluster false "集群"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *ClusterController) GetClusters() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	cluster := new(models.Cluster)
	json.Unmarshal(this.Ctx.Input.RequestBody, &cluster)
	this.Data["json"] = cluster.List(from, limit)
	this.ServeJSON(false)
}

// @Title UpdateCluster
// @Description Update Cluster
// @Param token header string true "authToken"
// @Param id path string "" true "Id"
// @Param body body models.Cluster true "集群"
// @Success 200 {object} models.Result
// @router /:id [put]
func (this *ClusterController) UpdateCluster() {
	id := this.GetString(":id")
	cluster := new(models.Cluster)
	json.Unmarshal(this.Ctx.Input.RequestBody, &cluster)
	cluster.Id = id
	this.Data["json"] = cluster.Update()
	this.ServeJSON(false)
}
