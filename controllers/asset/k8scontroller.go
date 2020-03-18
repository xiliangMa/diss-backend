package asset

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/models/k8s"
	"github.com/xiliangMa/diss-backend/utils"
)

// Asset K8S object api list
type K8SController struct {
	beego.Controller
}

// @Title GetClusters
// @Description Get Cluster List
// @Param token header string true "authToken"
// @Param user header string "admin" true "diss api 系统的登入用户"
// @Param body body k8s.Cluster false "集群"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /clusters [post]
func (this *K8SController) GetClusters() {
	accountName := models.Account_Admin
	user := this.Ctx.Input.Header("user")
	if user != models.Account_Admin && user != "" {
		accountUsers := models.AccountUsers{}
		accountUsers.UserName = user
		err, account := accountUsers.GetAccountByUser()
		accountName = account
		if err != nil {
			this.Data["json"] = models.Result{Code: utils.NoAccountUsersErr, Data: nil, Message: err.Error()}
			this.ServeJSON(false)
		}
	}

	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	cluster := new(k8s.Cluster)
	json.Unmarshal(this.Ctx.Input.RequestBody, &cluster)
	cluster.AccountName = accountName
	this.Data["json"] = cluster.List(from, limit)
	this.ServeJSON(false)

}

// @Title GetNameSpaceList
// @Description Get NameSpace List
// @Param token header string true "authToken"
// @Param user header string "admin" true "diss api 系统的登入用户"
// @Param clusterId path string "" true "clusterId"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /clusters/:clusterId/namespaces [post]
func (this *K8SController) GetNameSpaces() {
	user := this.Ctx.Input.Header("user")
	accountName := models.Account_Admin
	if user != models.Account_Admin && user != "" {
		accountUsers := models.AccountUsers{}
		accountUsers.UserName = user
		err, account := accountUsers.GetAccountByUser()
		accountName = account
		if err != nil {
			this.Data["json"] = models.Result{Code: utils.NoAccountUsersErr, Data: nil, Message: err.Error()}
			this.ServeJSON(false)
		}
	}

	clusterId := this.GetString(":clusterId")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	ns := new(k8s.NameSpace)
	ns.ClusterId = clusterId
	ns.AccountName = accountName
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
	this.Data["json"] = containerConfig.List(from, limit, false)
	this.ServeJSON(false)

}

// @Title GetContainerInfo
// @Description Get Container  Info
// @Param token header string true "authToken"
// @Param containerId path string "" true "containerId"
// @Success 200 {object} models.Result
// @router /containers/:containerId [post]
func (this *K8SController) GetContainerInfo() {
	containerId := this.GetString(":containerId")
	containerInfo := new(models.ContainerInfo)
	containerInfo.Id = containerId
	this.Data["json"] = containerInfo.List()
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

// @Title GetContainerPs
// @Description Get Container Ps  List
// @Param token header string true "authToken"
// @Param containerId path string "" true "containerId"
// @Param body body models.ContainerPs false "容器进程"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /containers/:containerId/containerps [post]
func (this *K8SController) GetContainerPs() {
	containerId := this.GetString(":containerId")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	containerPs := new(models.ContainerPs)
	json.Unmarshal(this.Ctx.Input.RequestBody, &containerPs)
	containerPs.ContainerId = containerId
	this.Data["json"] = containerPs.List(from, limit)
	this.ServeJSON(false)

}

// @Title ContainerImageInfo
// @Description Get ContainerImage Info  （根据当前 pod 下容器的 主机名 + 镜像名获取 镜像详细信息）
// @Param token header string true "authToken"
// @Param hostName path string "" true "hostName"
// @Param body body models.ImageInfo false "容器详细信息"
// @Success 200 {object} models.Result
// @router /:hostName/imageinfo [post]
func (this *K8SController) GetContainerImageInfo() {
	hostName := this.GetString(":hostName")

	imageInfo := new(models.ImageInfo)
	json.Unmarshal(this.Ctx.Input.RequestBody, &imageInfo)
	imageInfo.HostName = hostName

	this.Data["json"] = imageInfo.List()
	this.ServeJSON(false)

}
