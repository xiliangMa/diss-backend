package asset

import (
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
// @Param token header string true "auth token"
// @Param name query string "" false "name"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /clusters [post]
func (this *K8SController) GetClusters() {
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
// @Param token header string true "auth token"
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
// @Param token header string true "auth token"
// @Param nsName path string "" true "namespaceName"
// @Param podId path string "" true "podId"
// @Param name query string "" false "containerName"
// @Param imageName query string "" false "imageName"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /namespaces/:nsName/pods/:podId/containers [post]
func (this *K8SController) GetContainerConfig() {
	nsName := this.GetString(":nsName")
	podId := this.GetString(":podId")
	imageName := this.GetString("imageName")
	name := this.GetString("name")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	containerConfig := new(models.ContainerConfig)
	containerConfig.NameSpaceName = nsName
	containerConfig.PodId = podId
	containerConfig.Name = name
	containerConfig.ImageName = imageName
	this.Data["json"] = containerConfig.List(from, limit)
	this.ServeJSON(false)

}

// @Title GetContainerCmdHistorys
// @Description Get Container CmdHistory  List
// @Param token header string true "auth token"
// @Param containersId path string "" true "containersId"
// @Param command query string "" false "command"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /containers/:containersId/cmdhistorys [post]
func (this *K8SController) GetContainerCmdHistorys() {
	containersId := this.GetString(":containersId")
	command := this.GetString("command")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	cmdHistory := new(models.CmdHistory)
	cmdHistory.Command = command
	cmdHistory.ContainerId = containersId
	cmdHistory.Type = 1
	this.Data["json"] = cmdHistory.List(from, limit)
	this.ServeJSON(false)

}
