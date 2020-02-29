package asset

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/models/k8s"
)

// Asset K8S object api list
type K8SController struct {
	beego.Controller
}

// @Title GetClusters
// @Description Get Cluster List
// @Param token header string true "authToken"
// @Param body body k8s.Cluster false "集群"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /clusters [post]
func (this *K8SController) GetClusters() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	cluster := new(k8s.Cluster)
	json.Unmarshal(this.Ctx.Input.RequestBody, &cluster)
	this.Data["json"] = cluster.List(from, limit)
	this.ServeJSON(false)

}

// @Title GetNameSpaceList
// @Description Get NameSpace List
// @Param token header string true "authToken"
// @Param clusterId path string "" true "clusterId"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /clusters/:clusterId/namespaces [post]
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
// @Param token header string true "authToken"
// @Param nsName path string "" true "namespaceName"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /namespaces/:nsName/pods [post]
func (this *K8SController) GetPods() {
	nsName := this.GetString(":nsName")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	pod := new(k8s.Pod)
	pod.NameSpaceName = nsName
	this.Data["json"] = pod.List(from, limit)
	this.ServeJSON(false)

}

// @Title GetContainerList
// @Description Get pod Container List
// @Param token header string true "authToken"
// @Param nsName path string "" true "namespaceName"
// @Param podId path string "" true "podId"
// @Param body body models.ContainerConfig false "容器配置信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /namespaces/:nsName/pods/:podId/containers [post]
func (this *K8SController) GetContainerConfig() {
	nsName := this.GetString(":nsName")
	podId := this.GetString(":podId")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	containerConfig := new(models.ContainerConfig)
	json.Unmarshal(this.Ctx.Input.RequestBody, &containerConfig)
	containerConfig.NameSpaceName = nsName
	containerConfig.PodId = podId
	this.Data["json"] = containerConfig.List(from, limit)
	this.ServeJSON(false)

}

// @Title GetContainerCmdHistorys
// @Description Get Container CmdHistory  List
// @Param token header string true "authToken"
// @Param containerId path string "" true "containerId"
// @Param body body models.CmdHistory false "主机命令历史"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /containers/:containerId/cmdhistorys [post]
func (this *K8SController) GetContainerCmdHistorys() {
	containerId := this.GetString(":containerId")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	cmdHistory := new(models.CmdHistory)
	json.Unmarshal(this.Ctx.Input.RequestBody, &cmdHistory)
	cmdHistory.ContainerId = containerId
	cmdHistory.Type = 1
	this.Data["json"] = cmdHistory.List(from, limit)
	this.ServeJSON(false)

}

// @Title GetContainerTop
// @Description Get Container Top  List
// @Param token header string true "authToken"
// @Param containerId path string "" true "containerId"
// @Param body body models.ContainerTop false "容器进程"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /containers/:containerId/containertop [post]
func (this *K8SController) GetContainerTop() {
	containerId := this.GetString(":containerId")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	containerTop := new(models.ContainerTop)
	json.Unmarshal(this.Ctx.Input.RequestBody, &containerTop)
	containerTop.ContainerId = containerId
	this.Data["json"] = containerTop.List(from, limit)
	this.ServeJSON(false)

}
