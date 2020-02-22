package asset

import (
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/models/k8s"
	msl "github.com/xiliangMa/diss-backend/models/securitylog"
	ssl "github.com/xiliangMa/diss-backend/service/securitylog"
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
// @Param name query string "" false "name"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /:hostId/info [post]
func (this *HostController) GetHostInfoList() {
	id := this.GetString(":hostId")
	name := this.GetString("name")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	hostInfo := new(models.HostInfo)
	hostInfo.Id = id
	hostInfo.HostName = name
	this.Data["json"] = hostInfo.List(id, from, limit)
	this.ServeJSON(false)

}

// @Title GetHostPod
// @Description Get HostPod List
// @Param token header string true "auth token"
// @Param hostName path string "" true "hostName"
// @Param name query string "" false "podName"
// @Param from query int 0 false "from"
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
	pod.HostName = hostName
	this.Data["json"] = pod.List(from, limit)
	this.ServeJSON(false)

}

// @Title HostImage
// @Description Get HostImage List
// @Param token header string true "auth token"
// @Param hostId path string "" true "hostId"
// @Param name query string "" false "imageName"
// @Param from query int 0 false "from"
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
	imageConfig.HostId = hostId
	this.Data["json"] = imageConfig.List(from, limit)
	this.ServeJSON(false)

}

// @Title HostImageInfo
// @Description Get HostImage Info
// @Param token header string true "auth token"
// @Param hostId path string "" true "hostId"
// @Param imageId path string "" true "imageId"
// @Success 200 {object} models.Result
// @router /:hostId/images/:imageId [post]
func (this *HostController) GetHostImageInfo() {
	hostId := this.GetString(":hostId")
	imageId := this.GetString(":imageId")

	imageInfo := new(models.ImageInfo)
	imageInfo.HostId = hostId
	imageInfo.ImageId = imageId

	this.Data["json"] = imageInfo.List()
	this.ServeJSON(false)

}

// @Title HostContainerConfig
// @Description Get HostContainerConfig List
// @Param token header string true "auth token"
// @Param hostName path string "" true "hostName"
// @Param name query string "" false "containerName"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /:hostName/containers [post]
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
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /:hostId/ps [post]
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
// @Description Get HostContainer info
// @Param token header string true "auth token"
// @Param hostName path string "" true "hostName"
// @Param containerId path string "" true "containerId"
// @Success 200 {object} models.Result
// @router /:hostName/containers/:containerId [post]
func (this *HostController) GetHostContainerInfoList() {
	hostName := this.GetString(":hostName")
	containerId := this.GetString(":containerId")

	containerInfo := new(models.ContainerInfo)
	containerInfo.HostName = hostName
	containerInfo.Id = containerId
	this.Data["json"] = containerInfo.List()
	this.ServeJSON(false)

}

// @Title HostBenchMarkLog
// @Description Get HostBenchMarkLog List
// @Param token header string true "auth token"
// @Param hostId path string "" true "hostId"
// @Param bmtName query string "" false "bench mark template name"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /:hostId/hostbmls [post]
func (this *HostController) GetHostBenchMarkLogList() {
	bmtName := this.GetString("bmtName")
	hostId := this.GetString(":hostId")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	benchMarkLog := new(msl.BenchMarkLog)
	benchMarkLog.BenchMarkName = bmtName
	benchMarkLog.HostId = hostId
	this.Data["json"] = benchMarkLog.List(from, limit)
	this.ServeJSON(false)

}

// @Title HostBenchMarkLogInfo
// @Description Get HostBenchMarkLog Info
// @Param token header string true "auth token"
// @Param hostId path string "" true "hostId"
// @Param bmlId path string "" true "bench mark log id"
// @Param bmtName query string "" false "bench mark template name"
// @Success 200 {object} models.Result
// @router /:hostId/hostbmls/:bmlId [post]
func (this *HostController) GetHostBenchMarkLogInfo() {
	bmtName := this.GetString("bmtName")
	hostId := this.GetString(":hostId")
	bmlId := this.GetString(":bmlId")

	benchMarkLog := new(msl.BenchMarkLog)
	benchMarkLog.BenchMarkName = bmtName
	benchMarkLog.HostId = hostId
	benchMarkLog.Id = bmlId
	var securityLogService = ssl.SecurityLogService{benchMarkLog}
	this.Data["json"] = securityLogService.GetSecurityLogInfo()
	this.ServeJSON(false)

}
