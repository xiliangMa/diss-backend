package anchoreengine

import (
	"encoding/json"
	"github.com/beego/beego/v2/server/web"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/anchoreengine"
)

type AnchoreImageController struct {
	web.Controller
}

// @Title GetImageContent
// @Description Get Image Content
// @Param token header string true "authToken"
// @Param imageId path string "" true "imageId"
// @Param body body models.ImageSearchParams false "镜像信息"
// @Success 200 {object} models.Result
// @router /images/:imageId/content [post]
func (this *AnchoreImageController) GetImageContent() {
	imageId := this.GetString(":imageId")
	imageParams := new(models.ImageSearchParams)
	imageParams.ImageId = imageId
	json.Unmarshal(this.Ctx.Input.RequestBody, &imageParams)
	anchoreService := anchoreengine.AnchoreService{ImageParams: imageParams}
	this.Data["json"] = anchoreService.GetImageContent()
	this.ServeJSON(false)
}

// @Title GetImageVuln
// @Description Get Image Vuln
// @Param token header string true "authToken"
// @Param imageId path string "" true "imageId"
// @Param body body models.ImageSearchParams false "镜像信息"
// @Success 200 {object} models.Result
// @router /images/:imageId/vulns [post]
func (this *AnchoreImageController) GetImageVuln() {
	imageId := this.GetString(":imageId")
	imageParams := new(models.ImageSearchParams)
	imageParams.ImageId = imageId
	json.Unmarshal(this.Ctx.Input.RequestBody, &imageParams)
	anchoreService := anchoreengine.AnchoreService{ImageParams: imageParams}
	this.Data["json"] = anchoreService.GetImageVuln()
	this.ServeJSON(false)
}

// @Title GetImageVulnStatistics
// @Description Get Image Vuln Statistics
// @Param token header string true "authToken"
// @Param imageId path string "" true "imageId"
// @Param body body models.ImageSearchParams false "镜像信息"
// @Success 200 {object} models.Result
// @router /images/:imageId/vulnsstatistics [post]
func (this *AnchoreImageController) GetImageVulnStatistics() {
	imageId := this.GetString(":imageId")
	imageParams := new(models.ImageSearchParams)
	imageParams.ImageId = imageId
	json.Unmarshal(this.Ctx.Input.RequestBody, &imageParams)
	anchoreService := anchoreengine.AnchoreService{ImageParams: imageParams}
	this.Data["json"] = anchoreService.GetImageVulnStatistics()
	this.ServeJSON(false)
}

// @Title GetImageMetadata
// @Description Get Image Metadata
// @Param token header string true "authToken"
// @Param imageDigest path string "" true "imageDigest"
// @Param body body models.ImageSearchParams false "镜像信息"
// @Success 200 {object} models.Result
// @router /images/:imageDigest/metadata [post]
func (this *AnchoreImageController) GetImageMetadata() {
	imageDigest := this.GetString(":imageDigest")
	imageParams := new(models.ImageSearchParams)
	imageParams.ImageDigest = imageDigest
	json.Unmarshal(this.Ctx.Input.RequestBody, &imageParams)
	anchoreService := anchoreengine.AnchoreService{ImageParams: imageParams}
	this.Data["json"] = anchoreService.GetImageMetadata()
	this.ServeJSON(false)
}

// @Title GetImageSensitiveInfo
// @Description Get Image SensitiveInfo
// @Param token header string true "authToken"
// @Param imageId path string "" true "imageId"
// @Success 200 {object} models.Result
// @router /images/:imageId/sensitiveinfo [post]
func (this *AnchoreImageController) GetImageSensitiveInfo() {
	imageId := this.GetString(":imageId")
	imageParams := new(models.ImageSearchParams)
	imageParams.ImageId = imageId
	anchoreService := anchoreengine.AnchoreService{ImageParams: imageParams}
	this.Data["json"] = anchoreService.GetImageSensitiveInfo()
	this.ServeJSON(false)
}
