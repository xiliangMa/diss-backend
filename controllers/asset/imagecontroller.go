package asset

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
	"net/http"
)

// Asset image object api list
type ImageController struct {
	beego.Controller
}

// @Title GetContainers
// @Description Get Containers List
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
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

// @Title ImageDetail
// @Description Get ImageDetail Info
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param body body models.ImageDetail false "镜像信息"
// @Success 200 {object} models.Result
// @router /imagedetail [post]
func (this *ImageController) GetImageDetail() {
	imageDetail := new(models.ImageDetail)
	json.Unmarshal(this.Ctx.Input.RequestBody, &imageDetail)

	result := models.Result{Code: http.StatusOK}
	result.Data = imageDetail.Get()
	this.Data["json"] = result
	this.ServeJSON(false)

}
