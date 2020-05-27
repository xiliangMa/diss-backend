package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"os"
)

type LogConfig struct {
	Id            string `orm:"pk;" description:"(Log配置id)"`
	ConfigName    string `orm:"" description:"(配置项名称 支持的值 SysLogExport)"`
	Enabled       bool   `orm:"" description:"(是否启用)"`
	ServerUrl     string `orm:"" description:"(服务器url)"`
	ServerPort    string `orm:"" description:"(服务器端口)"`
	ExportedTypes string `orm:"" description:"(导出日志类型  多个枚举 以,分割)"`
}

type LogConfigInterface interface {
	Add() Result
	Update() Result
	Get() []*LogConfig
}

var GlobalLogConfig = map[string]*LogConfig{}

func (this *LogConfig) Get() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var logConfig []*LogConfig = nil
	var ResultData Result
	cond := orm.NewCondition()

	if this.ConfigName != "" {
		cond = cond.And("config_name", this.ConfigName)
	}

	_, err := o.QueryTable(utils.LogConfig).SetCond(cond).All(&logConfig)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetLogConfigErr
		logs.Error("Get LogConfig failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
	}

	//total, _ := o.QueryTable(utils.LogConfig).SetCond(cond).Count()
	//data := make(map[string]interface{})
	//data["item"] = logConfig
	//data["total"] = total
	//
	ResultData.Code = http.StatusOK
	ResultData.Data = logConfig
	return ResultData
}

// 此项目前为内部使用
func (this *LogConfig) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	uid, _ := uuid.NewV4()
	this.Id = uid.String()
	_, err := o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddLogConfigErr
		logs.Error("Add LogConfig failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *LogConfig) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditLogConfigErr
		logs.Error("Update LogConfig: %s failed, code: %d, err: %s", this.ConfigName, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func GetSyslogServerUrl() string {
	runMode := beego.AppConfig.String("RunMode")
	envRunMode := os.Getenv("RunMode")
	if envRunMode != "" {
		runMode = envRunMode
	}

	serverUrl := beego.AppConfig.String("syslog::SyslogServer")
	runMode = utils.Run_Mode_Prod
	if runMode == utils.Run_Mode_Prod {
		config, ok := GlobalLogConfig[Log_Config_SysLog_Export]
		if ok {
			serverUrl = config.ServerUrl + ":" + config.ServerPort
		}
	}
	return serverUrl
}
