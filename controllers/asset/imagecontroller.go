package asset

import (
	"encoding/json"
	"github.com/beego/beego/v2/server/web"
	"github.com/xiliangMa/diss-backend/models"
)

// Asset image object api list
type ImageController struct {
	web.Controller
}

// @Title GetContainers
// @Description Get Containers List
// @Param token header string true "authToken"
// @Param body body models.ContainerConfig false "容器配置信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /containers [post]
func (this *ImageController) GetContainersList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	containerConfig := new(models.ContainerConfig)
	json.Unmarshal(this.Ctx.Input.RequestBody, &containerConfig)
	this.Data["json"] = containerConfig.List(from, limit, false)
	this.ServeJSON(false)
}
