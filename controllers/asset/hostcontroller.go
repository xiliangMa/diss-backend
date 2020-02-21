package asset

import (
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/models/k8s"
)

// Asset host object api list
type HostController struct {
	beego.Controller
}

// host api list
// @Title GetHostConfig
// @Description Get HostConfig List
// @Param token header string true "auth token"
// @Param name query string "" false "name"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /hosts [post]
func (this *HostController) GetHostConfigList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	hostConfig := new(models.HostConfig)
	hostConfig.HostName = this.GetString("name")
	this.Data["json"] = hostConfig.List(from, limit)
	this.ServeJSON(false)

}

// @Title GetHostConfigInfo
// @Description Get HostConfigInfo
// @Param token header string true "auth token"
// @Param hostId path string "" true "hostId"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /hosts/:hostId/info [post]
func (this *HostController) GetHostInfoList() {
	id := this.GetString(":hostId")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	hostInfo := new(models.HostInfo)
	this.Data["json"] = hostInfo.List(id, from, limit)
	this.ServeJSON(false)

}

// @Title GetHostPod
// @Description Get HostPod List
// @Param token header string true "auth token"
// @Param hostName path string "" true "hostName"
// @Param name query string "" false "podName"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /hosts/:hostName/pods [post]
func (this *HostController) GetHostPodList() {
	name := this.GetString("name")
	hostName := this.GetString(":hostName")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	pod := new(k8s.Pod)
	pod.Name = name
	pod.HostName = hostName
	this.Data["json"] = pod.List(from, limit)
	this.ServeJSON(false)

}

// @Title HostImage
// @Description Get HostImage List
// @Param token header string true "auth token"
// @Param hostId path string "" true "hostId"
// @Param name query string "" false "imageName"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /hosts/:hostId/images [post]
func (this *HostController) GetHostImagesList() {
	name := this.GetString("name")
	hostId := this.GetString(":hostId")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	imageConfig := new(models.ImageConfig)
	imageConfig.Name = name
	this.Data["json"] = imageConfig.List(hostId, from, limit)
	this.ServeJSON(false)

}

// @Title HostContainerConfig
// @Description Get HostContainerConfig List
// @Param token header string true "auth token"
// @Param hostName path string "" true "hostName"
// @Param name query string "" false "containerName"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /hosts/:hostName/containers [post]
func (this *HostController) GetHostContainerConfigList() {
	name := this.GetString("name")
	hostName := this.GetString(":hostName")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	containerConfig := new(models.ContainerConfig)
	containerConfig.Name = name
	containerConfig.HostName = hostName
	this.Data["json"] = containerConfig.List(from, limit)
	this.ServeJSON(false)

}

// @Title HostPs
// @Description Get HostPs List
// @Param token header string true "auth token"
// @Param hostId path string "" true "hostId"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /hosts/:hostId/ps [post]
func (this *HostController) GetHostPsList() {
	hostId := this.GetString(":hostId")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	hostPs := new(models.HostPs)
	hostPs.HostId = hostId
	this.Data["json"] = hostPs.List(from, limit)
	this.ServeJSON(false)

}

// @Title HostContainerInfo
// @Description Get HostContainerInfo List
// @Param token header string true "auth token"
// @Param hostName path string "" true "hostName"
// @Param hostName path string "" true "containerId"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /hosts/:hostName/containers/:containerId [post]
func (this *HostController) GetHostContainerInfoList() {
	hostName := this.GetString(":hostName")
	containerId := this.GetString(":containerId")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	containerInfo := new(models.ContainerInfo)
	containerInfo.HostName = hostName
	containerInfo.Id = containerId
	this.Data["json"] = containerInfo.List(from, limit)
	this.ServeJSON(false)

}
