package system

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

type IntegrationController struct {
	beego.Controller
}

// @Title SystemIntegration Config
// @Description SystemIntegration Config
// @Param token header string true "authToken"
// @Param configName query string true "Config Name"
// @Param body body models.LogConfig false "日志配置信息"
// @Success 200 {object} models.Result
// @router /system/logconfig [post]
func (this *IntegrationController) AddLogConfig() {
	logConfig := new(models.LogConfig)
	configname := this.GetString("configName")
	json.Unmarshal(this.Ctx.Input.RequestBody, &logConfig)
	if configname != "" {
		logConfig.ConfigName = configname
	}

	this.Data["json"] = logConfig.Add()
	this.ServeJSON(false)
}

// @Title Update SystemIntegration Config
// @Description Update One SystemIntegration Config
// @Param token header string true "authToken"
// @Param body body models.LogConfig false "日志配置信息"
// @Success 200 {object} models.Result
// @router /system/logconfig [put]
func (this *IntegrationController) UpdateLogConfig() {
	logConfig := new(models.LogConfig)
	json.Unmarshal(this.Ctx.Input.RequestBody, &logConfig)

	result := logConfig.Update()
	this.Data["json"] = result
	logConfigObj := result.Data.(*models.LogConfig)
	models.GlobalLogConfig[models.Log_Config_SysLog_Export] = logConfigObj
	this.ServeJSON(false)
}

////内部Action，仅用于测试
// @Title Get SystemIntegration Config
// @Description Get One SystemIntegration Config [Inner内部操作]
// @Param token header string true "authToken"
// @Param configName query string true "Config Name"
// @Success 200 {object} models.Result
// @router /system/logconfig [get]
func (this *IntegrationController) InnerGetLogConfig() {
	logConfig := new(models.LogConfig)
	configName := this.GetString(":configName")
	logConfig.ConfigName = configName
	json.Unmarshal(this.Ctx.Input.RequestBody, &logConfig)

	fmt.Println("==============global log config============", models.GlobalLogConfig[models.Log_Config_SysLog_Export])
	this.Data["json"] = logConfig.InnerGet()
	this.ServeJSON(false)
}
