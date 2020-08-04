package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type SysConfig struct {
	Id    string `orm:"pk;" description:"(配置项id)"`
	Key   string `orm:"size(128)" description:"(配置项键)"`
	Value string `orm:"" description:"(配置项值 json)"`
}

type SysConfigInterface interface {
	Add() Result
	Update() Result
	List() Result
	Get() *SysConfig
}

func (this *SysConfig) Get() *SysConfig {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var sysConfig []*SysConfig = nil
	var ResultData Result
	cond := orm.NewCondition()

	if this.Key != "" {
		cond = cond.And("key", this.Key)
	}

	_, err := o.QueryTable(utils.SysConfig).SetCond(cond).All(&sysConfig)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetSysConfigErr
		logs.Error("Get SysConfig failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
	}
	if len(sysConfig) > 0 {
		return sysConfig[0]
	} else {
		return nil
	}
}

func (this *SysConfig) List() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var sysConfig []*SysConfig = nil
	var ResultData Result
	cond := orm.NewCondition()

	if this.Key != "" {
		cond = cond.And("key", this.Key)
	}

	_, err := o.QueryTable(utils.SysConfig).SetCond(cond).All(&sysConfig)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetSysConfigErr
		logs.Error("Get SysConfig failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
	}

	total, _ := o.QueryTable(utils.SysConfig).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = sysConfig
	if total == 0 {
		ResultData.Data = nil
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *SysConfig) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	uid, _ := uuid.NewV4()
	this.Id = uid.String()
	_, err := o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddSysConfigErr
		logs.Error("Add SysConfig failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *SysConfig) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditSysConfigErr
		logs.Error("Update SysConfig: %s failed, code: %d, err: %s", this.Key, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}
