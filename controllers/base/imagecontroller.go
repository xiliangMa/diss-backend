package base

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// 镜像接口列表
type ImageController struct {
	beego.Controller
}

// @Title GetImages
// @Description Get Images List
// @Param token header string true "authToken"
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
// @Param imageId path string "" true "imageId"
// @Success 200 {object} models.Result
// @router /:imageId [delete]
func (this *ImageController) DeleteImage() {
	imageId := this.GetString(":imageId")
	imageConfig := new(models.ImageConfig)
	imageConfig.Id = imageId
	this.Data["json"] = imageConfig.Delete()
	this.ServeJSON(false)
}

// @Title GetImageInfo
// @Description Get Image Info
// @Param token header string true "authToken"
// @Param imageId path string "" true "imageId"
// @Param body body models.ImageInfo false "镜像详细信息"
// @Success 200 {object} models.Result
// @router /:imageId/info [post]
func (this *ImageController) GetImageInfo() {
	imageId := this.GetString(":imageId")
	imageInfo := new(models.ImageInfo)
	json.Unmarshal(this.Ctx.Input.RequestBody, &imageInfo)
	imageInfo.ImageId = imageId
	this.Data["json"] = imageInfo.List()
	this.ServeJSON(false)
}

// @Title GetImageByDBType
// @Description Get Image Info
// @Param token header string true "authToken"
// @Param body body models.ImageConfig false "镜像配置信息"
// @Success 200 {object} models.Result
// @router /dbimage [post]
func (this *ImageController) GetImageByDBType() {
	imageConf := new(models.ImageConfig)
	json.Unmarshal(this.Ctx.Input.RequestBody, &imageConf)
	this.Data["json"] = imageConf.GetDBImageByType()
	this.ServeJSON(false)
}
