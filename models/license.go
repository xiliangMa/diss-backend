package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type LicenseFile struct {
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
	Id              string       `orm:"pk;" description:"(license module id)"`
	LicenseFile     *LicenseFile `orm:"rel(fk);null;" description:"(license file)"`
	ModuleCode      string       `orm:"" description:"(授权模块)"`
	LicenseCount    int          `orm:"" description:"(授权模块数量)"`
	LicenseExpireAt time.Time    `orm:"" description:"(授权结束时间)"`
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

type LicenseFileInterface interface {
	Add() Result
	Update() Result
	List(from, limit int) Result
}

func (licfile *LicenseFile) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	uuidlic,_ := uuid.NewV4()
	licfile.Id = uuidlic.String()
	_, err := o.Insert(licfile)

	licmodules := []*LicenseModule{}
	licmodules = licfile.LicenseModule
	licfile.LicenseModule = nil
	for _, licmodule := range licmodules{
		uuidmodule,_ := uuid.NewV4()
		licmodule.Id = uuidmodule.String()
		licmodule.LicenseFile = licfile
		o.Insert(licmodule)
	}
	//licfile.LicenseModule = licmodules

	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddLicenseFileErr
		logs.Error("Add LicenseFile failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = licfile
	return ResultData
}
