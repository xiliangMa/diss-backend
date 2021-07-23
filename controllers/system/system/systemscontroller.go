package system

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/nats"
	css "github.com/xiliangMa/diss-backend/service/system/system"
	"github.com/xiliangMa/diss-backend/utils"
	"io/ioutil"
	"net/http"
)

type SystemController struct {
	beego.Controller
}

// @Title UpLoadLogo
// @Description UpLoad Logo
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
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
		logoService.SaveDefaultLogo(fpath)
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

// @Title Restore Default Logo
// @Description Restore Default Logo
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Success 200 {object} models.Result
// @router /system/logo/restore [get]
func (this *SystemController) RestoreDefaultLogo() {
	logoService := new(css.LogoService)
	result := logoService.RestoreLogo()
	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title Check Logo isexist
// @Description Check kLogo IsExist
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
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
// @Param module header string true "moduleCode"
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
// @Param module header string true "moduleCode"
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
// @Param module header string true "moduleCode"
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

// @Title AddWarningWhiteList
// @Description Add WarningWhiteList Item
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param body body models.WarningWhiteList false "告警白名单"
// @Success 200 {object} models.Result
// @router /system/warnwhitelist [post]
func (this *SystemController) AddWarningWhiteList() {
	whiteListConfig := new(models.WarningWhiteList)
	json.Unmarshal(this.Ctx.Input.RequestBody, &whiteListConfig)

	this.Data["json"] = whiteListConfig.Add()
	this.ServeJSON(false)
}

// @Title UpdateWarningWhiteList
// @Description Update WarningWhiteList
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param id path string "" true "Id"
// @Param body body models.WarningWhiteList false "告警白名单"
// @Success 200 {object} models.Result
// @router /system/warnwhitelist/:id [put]
func (this *SystemController) UpdateWarningWhiteList() {
	id := this.GetString(":id")
	whiteListConfig := new(models.WarningWhiteList)
	json.Unmarshal(this.Ctx.Input.RequestBody, &whiteListConfig)
	whiteListConfig.Id = id
	result := whiteListConfig.Update()

	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title GetWarningWhiteList
// @Description Get WarningWhiteList
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param body body models.WarningWhiteList false "告警白名单"
// @Success 200 {object} models.Result
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @router /system/warnwhitelists [post]
func (this *SystemController) GetWarningWhiteList() {
	whiteListConfig := new(models.WarningWhiteList)
	json.Unmarshal(this.Ctx.Input.RequestBody, &whiteListConfig)

	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	this.Data["json"] = whiteListConfig.List(from, limit)
	this.ServeJSON(false)
}

// @Title DeleteWarningWhiteList
// @Description Delete WarningWhiteList
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param id path string "" true "id"
// @Success 200 {object} models.Result
// @router /system/warnwhitelist/:id [delete]
func (this *SystemController) DeleteWarningWhiteList() {
	id := this.GetString(":id")
	whiteListConfig := new(models.WarningWhiteList)
	whiteListConfig.Id = id

	this.Data["json"] = whiteListConfig.Delete()
	this.ServeJSON(false)
}

// @Title ImportWarnWhiteList
// @Description Import WarnInfo WhiteList
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param whitelistfile formData file true "whitelistfile"
// @Success 200 {object} models.Result
// @router /system/warnwhitelist/import [post]
func (this *SystemController) ImportWhiteList() {
	key := "whitelistfile"
	f, h, _ := this.GetFile(key)
	defer f.Close()

	whitelistService := new(css.WarnWhitelistService)
	result, fpath := whitelistService.Check(h)
	if result.Code != http.StatusOK {
		logs.Error("Upload Warninfo Whitelist fail, err: %s", result.Message)
	} else {
		err := this.SaveToFile(key, fpath)
		if err != nil {
			result.Code = utils.UploadWarnWhitelistErr
			result.Message = "UploadWarnWhitelistErr"
			logs.Error("Save Warninfo Whitelist fail, err: %s", err.Error())
		} else {
			// parse and save warninfo whitelist
			whitelistByte, err := ioutil.ReadAll(f)
			if err != nil {
				result.Code = utils.ImportWarnWhitelistErr
				result.Message = err.Error()
				logs.Error("Read Warninfo Whitelist file fail, err: %s", err)
			} else {
				whitelistService.WhitelistData = whitelistByte
				result = whitelistService.SaveList()
			}
		}
	}
	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title AddRuleDefine
// @Description Add RuleDefine Item
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param body body models.RuleDefine false "规则定义"
// @Success 200 {object} models.Result
// @router /system/ruledefine [post]
func (this *SystemController) AddRuleDefine() {
	ruleDefine := new(models.RuleDefine)
	json.Unmarshal(this.Ctx.Input.RequestBody, &ruleDefine)

	result := ruleDefine.Add()
	if result.Code == http.StatusOK {
		natsManager := models.Nats
		natsPubService := nats.NatsPubService{Conn: natsManager.Conn}
		natsPubService.Type = ruleDefine.Type
		natsPubService.RuleDefinePub()
	}

	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title UpdateRuleDefine
// @Description Update RuleDefine
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param id path string "" true "Id"
// @Param body body models.RuleDefine false "规则定义"
// @Success 200 {object} models.Result
// @router /system/ruledefine/:id [put]
func (this *SystemController) UpdateRuleDefine() {
	id, _ := this.GetInt64(":id")
	ruleDefine := new(models.RuleDefine)
	json.Unmarshal(this.Ctx.Input.RequestBody, &ruleDefine)
	ruleDefine.Id = id

	result := ruleDefine.Update()
	if result.Code == http.StatusOK {
		natsManager := models.Nats
		natsPubService := nats.NatsPubService{Conn: natsManager.Conn}
		natsPubService.Type = ruleDefine.Type
		natsPubService.RuleDefinePub()
	}

	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title GetRuleDefine
// @Description Get RuleDefine
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param body body models.RuleDefine false "规则定义"
// @Success 200 {object} models.Result
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @router /system/ruledefinelist [post]
func (this *SystemController) GetRuleDefine() {
	ruleDefine := new(models.RuleDefine)
	json.Unmarshal(this.Ctx.Input.RequestBody, &ruleDefine)

	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	this.Data["json"] = ruleDefine.List(from, limit)
	this.ServeJSON(false)
}

// @Title DeleteRuleDefine
// @Description Delete RuleDefine
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param id path string "" true "id"
// @Success 200 {object} models.RuleDefine
// @router /system/ruledefine/:id [delete]
func (this *SystemController) DeleteRuleDefine() {
	id, _ := this.GetInt64(":id")
	ruledefine := new(models.RuleDefine)
	ruledefine.Id = id

	this.Data["json"] = ruledefine.Delete()
	this.ServeJSON(false)
}

// @Title Version
// @Description version
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /version [post]
func (this *SystemController) Version() {
	from, _ := this.GetInt("from")
	limit, _ := this.GetInt("limit")

	version := new(models.Version)

	this.Data["json"] = version.List(from, limit)
	this.ServeJSON(false)

}

// @Title UpdateMetricData
// @Description Update Metric Data
// @Param token header string true "authToken"
// @Param hostIds path string "" false "hostId"
// @Param body body models.UpdateAssets true "更新配置设置"
// @Success 200 {object} models.Result
// @router /updatemetricdata [post]
func (this *SystemController) UpdateMetricData() {
	updateAssets := new(models.UpdateAssets)
	json.Unmarshal(this.Ctx.Input.RequestBody, &updateAssets)
	hostids := this.GetString(":hostIds")
	cmservice := css.ClientModuleService{}
	cmservice.ConfigForUpdate = updateAssets
	cmsdata := cmservice.GetUpdateAssetsConfig()

	natsManager := models.Nats
	natsPubService := nats.NatsPubService{Conn: natsManager.Conn}
	natsPubService.Message = cmsdata
	natsPubService.ClientSubject = hostids
	targetHosts := natsPubService.SendToManyHost()

	this.Data["json"] = targetHosts
	this.ServeJSON(false)

}
