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

// @Title GetHostConfig
// @Description Get HostConfig List
// @Param token header string true "auth token"
// @Param name query string "" false "name"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
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
// @router /:hostId/info [post]
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
// @router /:hostName/pods [post]
func (this *HostController) GetHostPodList() {
	name := this.GetString("name")
	hostName := this.GetString(":hostName")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	pod := new(k8s.Pod)
	pod.Name = name
	this.Data["json"] = pod.List(hostName, from, limit)
	this.ServeJSON(false)

}

// @Title HostImage
// @Description Get HostImage List
// @Param token header string true "auth token"
// @Param hostId path string "" true "hostId"
// @Param name query string "" false "imageName"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /:hostId/images [post]
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

// @Title HostContainer
// @Description Get HostContainer List
// @Param token header string true "auth token"
// @Param hostName path string "" true "hostName"
// @Param name query string "" false "containerName"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /:hostName/containers [post]
func (this *HostController) GetHostContainerList() {
	name := this.GetString("name")
	hostId := this.GetString(":hostName")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	containerConfig := new(models.ContainerConfig)
	containerConfig.Name = name
	this.Data["json"] = containerConfig.List(hostId, from, limit)
	this.ServeJSON(false)

}
