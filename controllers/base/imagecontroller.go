package base

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/base"
)

// 镜像接口列表
type ImageController struct {
	beego.Controller
}

// @Title GetImages
// @Description Get Images List
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param body body models.ImageConfig false "镜像配置信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *ImageController) GetImagesList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	imageConfig := new(models.ImageConfig)
	json.Unmarshal(this.Ctx.Input.RequestBody, &imageConfig)

	this.Data["json"] = imageConfig.List(from, limit)
	this.ServeJSON(false)
}

// @Title DeleteImage
// @Description Delete Image
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param imageId query string "" true "imageId"
// @Success 200 {object} models.Result
// @router / [delete]
func (this *ImageController) DeleteImage() {
	ids := this.GetString("imageId")

	imageDetail := new(models.ImageDetail)
	imageDetail.ImageConfigId = ids
	imageDetail.Delete()

	imageConfig := new(models.ImageConfig)
	imageConfig.Id = ids

	this.Data["json"] = imageConfig.Delete()
	this.ServeJSON(false)
}

// @Title GetImage
// @Description Get Image
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param body body models.ImageConfig false "镜像信息"
// @Success 200 {object} models.Result
// @router /get [post]
func (this *ImageController) GetImage() {
	imageConfig := new(models.ImageConfig)
	json.Unmarshal(this.Ctx.Input.RequestBody, &imageConfig)

	this.Data["json"] = imageConfig.Get()
	this.ServeJSON(false)
}

// @Title GetImageByDBType
// @Description Get Image Info
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param body body models.ImageConfig false "镜像配置信息"
// @Success 200 {object} models.Result
// @router /dbimage [post]
func (this *ImageController) GetImageByDBType() {
	imageConf := new(models.ImageConfig)
	json.Unmarshal(this.Ctx.Input.RequestBody, &imageConf)

	this.Data["json"] = imageConf.GetDBImageByType()
	this.ServeJSON(false)
}

// @Title BatchImportImage
// @Description Batch Import Image
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param body body models.ImageConfig false "镜像配置信息"
// @Success 200 {object} models.Result
// @router /import [post]
func (this *ImageController) BatchImportImage() {
	imageConfig := new(models.ImageConfig)
	json.Unmarshal(this.Ctx.Input.RequestBody, &imageConfig)

	ics := base.ImageConfigService{ImageConfig: imageConfig}

	this.Data["json"] = ics.BatchImportImage()
	this.ServeJSON(false)
}

// @Title GetNamespaces
// @Description GetNamespaces
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param body body models.ImageConfig false "镜像配置信息"
// @Success 200 {object} models.Result
// @router /namespaces [post]
func (this *ImageController) GetNamespaces() {
	imageConfig := new(models.ImageConfig)
	json.Unmarshal(this.Ctx.Input.RequestBody, &imageConfig)

	ics := base.ImageConfigService{ImageConfig: imageConfig}

	this.Data["json"] = ics.GetNamespaces()
	this.ServeJSON(false)
}
