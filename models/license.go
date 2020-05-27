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
	Id           string           `orm:"pk;" description:"(序列号)"`
	ProductName  string           `orm:"" description:"(产品名称)"`
	CustomerName string           `orm:"" description:"(许可对象)"`
	Type         int              `orm:"" description:"(授权类型 0测试 1正式)"`
	BuyAt        time.Time        `orm:"null;type(datetime)" description:"(授权购买时间)"`
	ActiveAt     time.Time        `orm:"null;auto_now;type(datetime)" description:"(激活时间)"`
	Modules      []*LicenseModule `orm:"reverse(many);null" description:"(授权的模块)"`
}

type LicenseHistory struct {
	Id          string    `orm:"pk;" description:"(history id)"`
	LicenseJson string    `orm:"" description:"(license json 文件)"`
	UpdateTime  time.Time `orm:"null;auto_now;type(datetime)" description:"(更新时间)"`
}

type LinceseModuleInterface interface {
	Add()
	Remove(string)
}

type LicenseConfigInterface interface {
	Add() Result
	Update() Result
	Get() Result
}

type LicenseModule struct {
	Id              string         `orm:"pk;" description:"(license module id)"`
	LicenseConfig   *LicenseConfig `orm:"rel(fk);null" description:"(license file)"`
	ModuleCode      string         `orm:"" description:"(授权模块)"`
	LicenseCount    int            `orm:"" description:"(授权模块数量)"`
	LicenseExpireAt time.Time      `orm:"" description:"(授权结束时间)"`
}

type LicenseHistoryInterface interface {
	Add() Result
	List(from, limit int) Result
}

func (this *LicenseModule)  Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	_, err := o.Insert(this)

	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.ImportLicenseFileErr
		logs.Error("Import License failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *LicenseModule)  Remove(licid string) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()
	cond = cond.And("license_config_id", licid)
	_, err := o.QueryTable(utils.LicenseModule).SetCond(cond).Delete()

	if err != nil  {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteLicenseModuleErr
		logs.Error("Import License failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}


func (this *LicenseConfig) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	err := o.Begin()
	licmodules := this.Modules
	for _, licmodule := range licmodules {
		uuidmodule, _ := uuid.NewV4()
		tmplicmodule := LicenseModule{Id:uuidmodule.String(), LicenseConfig:this, LicenseCount:licmodule.LicenseCount,LicenseExpireAt:licmodule.LicenseExpireAt,ModuleCode:licmodule.ModuleCode}
		_ = tmplicmodule.Add()
		if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
			o.Rollback()
		}
	}

	_, err = o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		o.Rollback()
	}
	err = o.Commit()

	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.ImportLicenseFileErr
		logs.Error("Import License failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *LicenseConfig) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	err := o.Begin()

	_, err = o.Update(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		o.Rollback()
	}

	licmodules := this.Modules
	tmpmodule := LicenseModule{}
	tmpmodule.Remove(this.Id)
	for _, licmodule := range licmodules {
		uuidmodule, _ := uuid.NewV4()
		tmplicmodule := LicenseModule{Id:uuidmodule.String(), LicenseConfig:this, LicenseCount:licmodule.LicenseCount,LicenseExpireAt:licmodule.LicenseExpireAt, ModuleCode:licmodule.ModuleCode}
		_ = tmplicmodule.Add()
		if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
			o.Rollback()
		}
	}

	err = o.Commit()

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

func (this *LicenseConfig) Get() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	var logConfigData []*LicenseConfig = nil

	cond := orm.NewCondition()

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}

	_, err := o.QueryTable(utils.LicenseConfig).SetCond(cond).RelatedSel().All(&logConfigData)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetLogConfigErr
		logs.Error("Get license failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
	}
	for _, logCofing := range logConfigData {
		o.LoadRelated(logCofing, "LicenseModule")
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = logConfigData
	return ResultData
}

func (this *LicenseHistory) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	uuidlic, _ := uuid.NewV4()
	this.Id = uuidlic.String()
	_, err := o.Insert(this)

	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddLicenseHistoryErr
		logs.Error("Add License history, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *LicenseHistory) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var licenseHistoryList []*LicenseHistory = nil
	var ResultData Result
	var err error
	cond := orm.NewCondition()
	if this.LicenseJson != "" {
		cond = cond.And("name__icontains", this.LicenseJson)
	}
	_, err = o.QueryTable(utils.LicenseHistory).SetCond(cond).Limit(limit, from).All(&licenseHistoryList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetLicenseHistoryErr
		logs.Error("Get License History List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.LicenseHistory).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = licenseHistoryList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}
