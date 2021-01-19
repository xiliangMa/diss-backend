package base

import (
	"encoding/json"
	"github.com/beego/beego/v2/server/web"
	"github.com/xiliangMa/diss-backend/models"
)

// 容器接口列表
type ContainerController struct {
	web.Controller
}

// @Title GetContainers
// @Description Get Containers List
// @Param token header string true "authToken"
// @Param body body models.ContainerConfig false "容器配置信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *ContainerController) GetContainersList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	containerConfig := new(models.ContainerConfig)
	json.Unmarshal(this.Ctx.Input.RequestBody, &containerConfig)
	this.Data["json"] = containerConfig.List(from, limit, false)
	this.ServeJSON(false)
}

// @Title DeleteContainer
// @Description Delete Container
// @Param token header string true "authToken"
// @Param containerId path string "" true "containerId"
// @Success 200 {object} models.Result
// @router /:containerId [delete]
func (this *ContainerController) DeleteContainer() {
	containerId := this.GetString(":containerId")
	containerConfig := new(models.ContainerConfig)
	//json.Unmarshal(this.Ctx.Input.RequestBody, &containerConfig)
	containerConfig.Id = containerId
	this.Data["json"] = containerConfig.Delete()
	this.ServeJSON(false)
}

// @Title GetContainerInfo
// @Description Get Container Info
// @Param token header string true "authToken"
// @Param containerId path string "" true "containerId"
// @Param body body models.ContainerInfo false "容器详细信息"
// @Success 200 {object} models.Result
// @router /:containerId/info [post]
func (this *ContainerController) GetContainerInfo() {
	containerId := this.GetString(":containerId")
	containerInfo := new(models.ContainerInfo)
	json.Unmarshal(this.Ctx.Input.RequestBody, &containerInfo)
	containerInfo.Id = containerId
	this.Data["json"] = containerInfo.List()
	this.ServeJSON(false)
}
