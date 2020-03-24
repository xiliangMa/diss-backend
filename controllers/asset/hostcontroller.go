package asset

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/models/k8s"
	msl "github.com/xiliangMa/diss-backend/models/securitylog"
	"github.com/xiliangMa/diss-backend/utils"
)

// Asset host object api list
type HostController struct {
	beego.Controller
}

// @Title GetHostConfig
// @Description Get HostConfig List
// @Param token header string true "authToken"
// @Param user query string "admin" true "diss api 系统的登入用户 如果用户all, 直接根据租户查询"
// @Param body body models.HostConfig false "主机配置信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *HostController) GetHostConfigList() {
	accountName := models.Account_Admin
	user := this.GetString("user")
	if user != models.Account_Admin && user != "" {
		accountUsers := new(models.AccountUsers)
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
	hostConfig := new(models.HostConfig)
	hostConfig.AccountName = accountName
	json.Unmarshal(this.Ctx.Input.RequestBody, &hostConfig)
	this.Data["json"] = hostConfig.List(from, limit)
	this.ServeJSON(false)

}

// @Title GetHostInfo
// @Description Get HostInfo
// @Param token header string true "authToken"
// @Param hostId path string "" true "hostId"
// @Success 200 {object} models.Result
// @router /:hostId [post]
func (this *HostController) GetHostInfoList() {
	id := this.GetString(":hostId")
	hostInfo := new(models.HostInfo)
	hostInfo.Id = id
	this.Data["json"] = hostInfo.List()
	this.ServeJSON(false)

}

// @Title GetHostPod
// @Description Get HostPod List
// @Param token header string true "authToken"
// @Param hostName path string "" true "主机名"
// @Param body body k8s.Pod false "Pod 信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /:hostName/pods [post]
func (this *HostController) GetHostPodList() {
	hostName := this.GetString(":hostName")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	pod := new(k8s.Pod)
	json.Unmarshal(this.Ctx.Input.RequestBody, &pod)
	pod.HostName = hostName
	this.Data["json"] = pod.List(from, limit)
	this.ServeJSON(false)

}

// @Title HostImage
// @Description Get HostImage List
// @Param token header string true "authToken"
// @Param hostId path string "" true "hostId"
// @Param body body models.ImageConfig false "镜像配置信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /:hostId/images [post]
func (this *HostController) GetHostImagesList() {
	hostId := this.GetString(":hostId")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	imageConfig := new(models.ImageConfig)
	json.Unmarshal(this.Ctx.Input.RequestBody, &imageConfig)
	imageConfig.HostId = hostId
	this.Data["json"] = imageConfig.List(from, limit)
	this.ServeJSON(false)

}

// @Title HostCmdHistory
// @Description Get HostCmdHistory List
// @Param token header string true "authToken"
// @Param hostId path string "" true "hostId"
// @Param body body models.CmdHistory false "主机命令历史信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /:hostId/cmdhistory [post]
func (this *HostController) GetHostCmdHistoryList() {
	hostId := this.GetString(":hostId")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	cmdHistory := new(models.CmdHistory)
	json.Unmarshal(this.Ctx.Input.RequestBody, &cmdHistory)
	cmdHistory.HostId = hostId
	cmdHistory.Type = 0
	this.Data["json"] = cmdHistory.List(from, limit)
	this.ServeJSON(false)

}

// @Title HostImageInfo
// @Description Get HostImage Info
// @Param token header string true "authToken"
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
// @Param token header string true "authToken"
// @Param hostName path string "" true "hostName"
// @Param body body models.ContainerConfig false "容器配置信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /:hostName/containers [post]
func (this *HostController) GetHostContainerConfigList() {
	hostName := this.GetString(":hostName")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	containerConfig := new(models.ContainerConfig)
	json.Unmarshal(this.Ctx.Input.RequestBody, &containerConfig)
	containerConfig.HostName = hostName
	this.Data["json"] = containerConfig.List(from, limit, false)
	this.ServeJSON(false)

}

// @Title HostPs
// @Description Get HostPs List
// @Param token header string true "authToken"
// @Param hostId path string "" true "hostId"
// @Param body body models.HostPs false "主机进程"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /:hostId/ps [post]
func (this *HostController) GetHostPsList() {
	hostId := this.GetString(":hostId")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	hostPs := new(models.HostPs)
	json.Unmarshal(this.Ctx.Input.RequestBody, &hostPs)
	hostPs.HostId = hostId
	this.Data["json"] = hostPs.List(from, limit)
	this.ServeJSON(false)

}

// @Title HostContainerInfo
// @Description Get HostContainer info
// @Param token header string true "authToken"
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
// @Param token header string true "authToken"
// @Param hostId path string "" true "hostId"
// @Param body body securitylog.BenchMarkLog false "基线日志信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /:hostId/hostbmls [post]
func (this *HostController) GetHostBenchMarkLogList() {
	hostId := this.GetString(":hostId")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	benchMarkLog := new(msl.BenchMarkLog)
	json.Unmarshal(this.Ctx.Input.RequestBody, &benchMarkLog)
	benchMarkLog.HostId = hostId
	this.Data["json"] = benchMarkLog.List(from, limit)
	this.ServeJSON(false)

}

// @Title HostBenchMarkLogInfo
// @Description Get HostBenchMarkLog Info
// @Param token header string true "authToken"
// @Param hostId path string "" true "hostId"
// @Param bmlId path string "" true "benchMarkLogId"
// @Success 200 {object} models.Result
// @router /:hostId/hostbmls/:bmlId [post]
func (this *HostController) GetHostBenchMarkLogInfo() {
	hostId := this.GetString(":hostId")
	bmlId := this.GetString(":bmlId")

	benchMarkLog := new(msl.BenchMarkLog)
	benchMarkLog.HostId = hostId
	benchMarkLog.Id = bmlId
	//var securityLogService = ssl.SecurityLogService{benchMarkLog, nil}
	//this.Data["json"] = securityLogService.GetHostBenchMarkLogInfo()
	this.Data["json"] = benchMarkLog.List(0, 0)
	this.ServeJSON(false)

}
