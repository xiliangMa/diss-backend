package base

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// 镜像阻断列表
type ImageBlockingController struct {
	beego.Controller
}

// @Title GetImageBlockingList
// @Description Get Image Blocking List
// @Param token header string true "authToken"
// @Param body body models.ImageBlocking false "镜像配置信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *ImageBlockingController) GetImageBlockingList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	imageBlocking := new(models.ImageBlocking)
	json.Unmarshal(this.Ctx.Input.RequestBody, &imageBlocking)

	this.Data["json"] = imageBlocking.List(from, limit)
	this.ServeJSON(false)
}

// @Title ImageBlocking operation
// @Description Image Blocking operation
// @Param token header string true "authToken"
// @Param body body models.ImageBlocking true "ImageBlocking"
// @Success 200 {object} models.Result
// @router /operation [post]
func (this *ImageBlockingController) Operation() {

	ib := new(models.ImageBlocking)
	json.Unmarshal(this.Ctx.Input.RequestBody, &ib)

	this.Data["json"] = ib.Add()
	this.ServeJSON(false)
}

// @Title Get ImageBlocking
// @Description Get Image Blocking
// @Param token header string true "authToken"
// @Param imageId path string "" true "imageId"
// @Success 200 {object} models.Result
// @router /:imageId [post]
func (this *ImageBlockingController) Get() {
	imageId := this.GetString(":imageId")
	ib := new(models.ImageBlocking)

	ib.ImageId = imageId

	this.Data["json"] = ib.Get()
	this.ServeJSON(false)
}
