package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type Registry struct {
	Id          int64  `orm:"pk;auto" description:"(仓库id)"`
	AccountName string `orm:"" description:"(租户)"`
	Name        string `orm:"size(64)" description:"(仓库名)"`
	Description string `orm:"size(256)" description:"(描述/备注)"`
	Type        string `orm:"size(32)" description:"(仓库类型)"`
	Url         string `orm:"size(512)" dqescription:"(地址)"`
	User        string `orm:"size(32)" description:"(用户名)"`
	Pwd         string `orm:"size(128)" description:"(密码)"`
	Insecure    bool   `orm:"default(true)" description:"(验证远程证书)"`
	CreateTime  int64  `orm:"default(0)" description:"(创建时间)"`
	//ImageConfig []*ImageConfig `orm:"reverse(many);default(null)" description:"(镜像)"`
}

type RegistryInterface interface {
	Add() Result
	Get() Result
	Delete() Result
	List(from, limit int) Result
}

func (this *Registry) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	var err error

	if r := this.Get(); r.Data != nil {
		ResultData.Message = "Name is already"
		ResultData.Code = utils.GetRegistryErr
		logs.Error("Get Registry failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	if this.CreateTime == 0 {
		this.CreateTime = time.Now().UnixNano()
	}

	_, err = o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddRegistryErr
		logs.Error("Add Registry failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Registry) Get() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	obj := new(Registry)
	cond := orm.NewCondition()

	if this.Id != 0 {
		cond = cond.And("id", this.Id)
	}

	if this.Name != "" {
		cond = cond.And("name", this.Name)
	}

	err := o.QueryTable(utils.Registry).SetCond(cond).RelatedSel().One(obj)

	if err != nil {
		ResultData.Code = utils.GetRegistryErr
		ResultData.Message = err.Error()
		return ResultData
	}

	data := make(map[string]interface{})
	data[Result_Items] = obj

	ResultData.Code = http.StatusOK
	ResultData.Data = data

	return ResultData
}

func (this *Registry) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var registryList []*Registry
	var ResultData Result
	var err error

	cond := orm.NewCondition()

	if this.AccountName != "" && this.AccountName != Account_Admin {
		cond = cond.And("account_name", this.AccountName)
	}

	if this.Name != "" {
		cond = cond.And("name__contains", this.Name)
	}

	_, err = o.QueryTable(utils.Registry).RelatedSel().SetCond(cond).Limit(limit, from).All(&registryList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetRegistryErr
		logs.Error("Get Registry List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.Registry).SetCond(cond).Count()
	data := make(map[string]interface{})
	data[Result_Total] = total
	data[Result_Items] = registryList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *Registry) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()

	if this.Id != 0 {
		cond = cond.And("id", this.Id)
	}
	_, err := o.QueryTable(utils.Registry).SetCond(cond).Delete()
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteRegistryErr
		logs.Error("Delete Registry failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}
