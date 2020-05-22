package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type LicenseConfig struct {
	Id              string           `orm:"pk;" description:"(license file id)"`
	ProductName     string           `orm:"" description:"(产品名称)"`
	CustomerName    string           `orm:"" description:"(许可对象)"`
	LicenseType     int              `orm:"" description:"(授权类型 0测试 1正式)"`
	LicenseUuid     string           `orm:"" description:"(序列号)"`
	LicenseBuyAt    time.Time        `orm:"null;type(datetime)" description:"(授权购买时间)"`
	LicenseActiveAt time.Time        `orm:"null;auto_now;type(datetime)" description:"(激活时间)"`
	LicenseModule   []*LicenseModule `orm:"reverse(many);null" description:"(授权的模块)"`
}

type LicenseModule struct {
	Id              string         `orm:"pk;" description:"(license module id)"`
	LicenseFile     *LicenseConfig `orm:"rel(fk);null;" description:"(license file)"`
	ModuleCode      string         `orm:"" description:"(授权模块)"`
	LicenseCount    int            `orm:"" description:"(授权模块数量)"`
	LicenseExpireAt time.Time      `orm:"" description:"(授权结束时间)"`
}

var LicenseModuleCodeMap = map[string]string{
	"ImageScan":           "镜像仓库扫描",
	"DockerBenchMark":     "Docker基线扫描",
	"KubernetesBenchMark": "K8s基线扫描",
	"IntrudeDetectScan":   "入侵扫描",
	"SecurityAudit":       "安全审计",
	"DockerVirusScan":     "Docker病毒扫描",
	"HostVirusScan":       "主机病毒扫描",
	"SC_LeakScan":         "漏洞扫描",
}

type LicenseConfigInterface interface {
	Add() Result
	Update() Result
	Get() Result
	List(from, limit int) Result
}

func (this *LicenseConfig) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	uuidlic, _ := uuid.NewV4()
	this.Id = uuidlic.String()
	_, err := o.Insert(this)

	licmodules := []*LicenseModule{}
	licmodules = this.LicenseModule
	this.LicenseModule = nil
	for _, licmodule := range licmodules {
		uuidmodule, _ := uuid.NewV4()
		licmodule.Id = uuidmodule.String()
		licmodule.LicenseFile = this
		o.Insert(licmodule)
	}

	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.ImportLicenseFileErr
		logs.Error("Import LicenseFile failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *LicenseConfig) Get() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	var logConfigData []*LicenseConfig = nil

	cond := orm.NewCondition()

	if this.LicenseUuid != "" {
		cond = cond.And("license_uuid", this.LicenseUuid)
	}

	err := o.QueryTable(utils.LicenseConfig).SetCond(cond).RelatedSel().One(&logConfigData)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetLogConfigErr
		logs.Error("Get LogConfig failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
	}
	for _, logCofing := range logConfigData {
		o.LoadRelated(logCofing, "LicenseModule")
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = logConfigData
	return ResultData
}

func (this *LicenseConfig) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditLicenseErr
		logs.Error("Update license: %s failed, code: %d, err: %s", this.Id, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}
