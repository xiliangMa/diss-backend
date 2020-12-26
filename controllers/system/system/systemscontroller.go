package system

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	css "github.com/xiliangMa/diss-backend/service/system/system"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type SystemController struct {
	beego.Controller
}

// @Title UpLoadLogo
// @Description UpLoad Logo
// @Param token header string true "authToken"
// @Param logo formData file true "logo"
// @Success 200 {object} models.Result
// @router /system/logo [post]
func (this *SystemController) UploadLogo() {
	key := "logo"
	f, h, _ := this.GetFile(key)
	defer f.Close()
	logoService := new(css.LogoService)
	result, fpath := logoService.Check(h)
	if result.Code != http.StatusOK {
		logs.Error("Upload logo  fail, err: %s", result.Message)
	} else {
		err := this.SaveToFile(key, fpath)
		if err != nil {
			result.Code = utils.UploadLogoErr
			result.Message = "UploadLogoErr"
			logs.Error("Upload logo  fail, err: %s", err.Error())
		} else {
			result.Code = http.StatusOK
		}
	}
	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title Check Logo isexist
// @Description Check kLogo IsExist
// @Param token header string true "authToken"
// @Success 200 {object} models.Result
// @router /system/logo/isexist [get]
func (this *SystemController) CheckLogoIsExist() {
	logoService := new(css.LogoService)
	result := logoService.CheckLogoIsExist()
	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title AddSysConfig
// @Description Add System Config
// @Param token header string true "authToken"
// @Param body body models.SysConfig false "系统配置信息"
// @Success 200 {object} models.Result
// @router /system/sysconfig [post]
func (this *SystemController) AddSysConfig() {
	sysConfig := new(models.SysConfig)
	json.Unmarshal(this.Ctx.Input.RequestBody, &sysConfig)

	this.Data["json"] = sysConfig.Add()
	this.ServeJSON(false)
}

// @Title GetSysConfigs
// @Description Get System Config
// @Param token header string true "authToken"
// @Param key query string false "配置项键名"
// @Success 200 {object} models.Result
// @router /system/sysconfig [get]
func (this *SystemController) GetSysConfigs() {
	sysConfig := new(models.SysConfig)
	configName := this.GetString("key")
	sysConfig.Key = configName
	json.Unmarshal(this.Ctx.Input.RequestBody, &sysConfig)

	this.Data["json"] = sysConfig.List()
	this.ServeJSON(false)
}

// @Title UpdateSysConfig
// @Description Update System Config
// @Param token header string true "authToken"
// @Param id path string "" true "Id"
// @Param body body models.SysConfig false "系统配置信息"
// @Success 200 {object} models.Result
// @router /system/sysconfig/:id [put]
func (this *SystemController) UpdateSysConfig() {
	id := this.GetString(":id")
	sysConfig := new(models.SysConfig)
	json.Unmarshal(this.Ctx.Input.RequestBody, &sysConfig)
	sysConfig.Id = id
	result := sysConfig.Update()

	sysConfigService := css.SysConfigService{}
	sysConfigService.SysConfig = sysConfig
	sysConfigService.RefreshConnections()

	this.Data["json"] = result
	this.ServeJSON(false)
}
