package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type NetworkPolicy struct {
	Id            string `orm:"pk;" description:"(id)"`
	Name          string `orm:"size(128)" description:"(名)"`
	AccountName   string `orm:"size(32)" description:"(租户)"`
	ClusterName   string `orm:"size(32)" description:"(集群名)"`
	NameSpaceName string `orm:"size(255);default(null);" description:"(命名空间)"`
	KMetaData     string `orm:"" description:"(源数据)"`
	KSpec         string `orm:"" description:"(Spec数据)"`
}

type NetworkPolicyInterface interface {
	Add() Result
	Delete() Result
	Update() Result
	List(from, limit int) Result
}

func (this *NetworkPolicy) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	var err error
	_, err = o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddNetworkPolicyErr
		logs.Error("Add NetworkPolicy failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *NetworkPolicy) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var NetworkPolicyList []*NetworkPolicy = nil
	var ResultData Result
	var err error
	cond := orm.NewCondition()
	if this.Name != "" {
		cond = cond.And("name__icontains", this.Name)
	}
	_, err = o.QueryTable(utils.NetworkPolicy).SetCond(cond).Limit(limit, from).All(&NetworkPolicyList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetNetworkPolicyErr
		logs.Error("Get NetworkPolicy List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.NetworkPolicy).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = NetworkPolicyList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *NetworkPolicy) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditNetworkPolicyErr
		logs.Error("Update NetworkPolicy: %s failed, code: %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *NetworkPolicy) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if this.ClusterName != "" {
		cond = cond.And("cluster_name", this.ClusterName)
	}
	_, err := o.QueryTable(utils.NetworkPolicy).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteNetworkPolicyErr
		logs.Error("Delete NetworkPolicy failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}
