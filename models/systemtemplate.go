package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type SystemTemplate struct {
	Id          string `orm:"pk;" description:"(基线id)"`
	Account     string `orm:"default(admin)" description:"(租户)"`
	Name        string `orm:"" description:"(名称)"`
	Description string `orm:"" description:"(描述)"`
	Type        string `orm:"" description:"(类型)"`
	Version     string `orm:"null" description:"(版本)"`
	Commands    string `orm:"null;" description:"(操作命令)"`
	Status      int    `orm:"default(1);" description:"(类型 停用 0  启用 1)"`
	IsDefault   bool   `orm:"default(false);" description:"(默认系统策略)"`
}

type SystemTemplateInterface interface {
	Add() Result
	List() Result
	Delete() Result
	Update() Result
	GetDefaultTemplate() map[string]*SystemTemplate
}

func (this *SystemTemplate) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddSYSTemplateErr
		logs.Error("Add SystemTemplate failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *SystemTemplate) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var systemTemplateList []*SystemTemplate
	var ResultData Result
	var err error
	cond := orm.NewCondition()

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if this.Account != "" {
		cond = cond.And("account", this.Account)
	}
	if this.Name != "" {
		cond = cond.And("name__contains", this.Name)
	}
	if this.Type != "" {
		cond = cond.And("type", this.Type)
	}
	if this.Version != All {
		cond = cond.And("name__contains", this.Name)
	}
	if this.Commands != "" {
		cond = cond.And("commands__contains", this.Commands)
	}
	if this.Status != TMP_Status_ALl {
		cond = cond.And("status", this.Status)
	}
	_, err = o.QueryTable(utils.SYSTemplate).SetCond(cond).RelatedSel().Limit(limit, from).All(&systemTemplateList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetSYSTemplateErr
		logs.Error("Get SystemTemplate List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.SYSTemplate).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = systemTemplateList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *SystemTemplate) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	_, err := o.QueryTable(utils.SYSTemplate).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteSYSTemplateErr
		logs.Error("Delete SYSTemplateErr failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}

func (this *SystemTemplate) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditSYSTemplateErr
		logs.Error("Update SYSTemplateErr: %s failed, code: %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *SystemTemplate) GetDefaultTemplate() map[string]*SystemTemplate {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var systemTemplateList []*SystemTemplate
	defaultTemplateList := make(map[string]*SystemTemplate)
	var err error
	cond := orm.NewCondition()
	cond = cond.And("is_default", true)
	_, err = o.QueryTable(utils.SYSTemplate).SetCond(cond).RelatedSel().All(&systemTemplateList)
	if err != nil {
		logs.Error("Get SystemTemplate List failed, code: %d, err: %s", utils.GetSYSTemplateErr, err.Error())
		return nil
	}

	for _, systemTemplate := range systemTemplateList {
		defaultTemplateList[systemTemplate.Type] = systemTemplate
	}
	return defaultTemplateList
}
