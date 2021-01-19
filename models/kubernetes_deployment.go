package models

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type Deployment struct {
	Id            string `orm:"pk;" description:"(id)"`
	Name          string `orm:"size(128)" description:"(Deployment名)"`
	AccountName   string `orm:"size(32)" description:"(租户)"`
	NameSpaceName string `orm:"size(255);default(null);" description:"(命名空间)"`
	HostName      string `orm:"size(64);default(null);" description:"(主机名)"`
	ClusterName   string `orm:"size(32)" description:"(集群名)"`
	KMetaData     string `orm:"" description:"(源数据)"`
	KSpec         string `orm:"" description:"(Spec数据)"`
	KStatus       string `orm:"" description:"(状态数据)"`
}

type DeploymentInterface interface {
	Add() Result
	Delete() Result
	Update() Result
	List(from, limit int) Result
}

func (this *Deployment) Add() Result {
	o := orm.NewOrm()
	var ResultData Result
	var DeploymentList []*Deployment
	var err error
	cond := orm.NewCondition()
	if this.Name != "" {
		cond = cond.And("id", this.Id)
	}

	_, err = o.QueryTable(utils.Deployment).SetCond(cond).All(&DeploymentList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetDeploymentErr
		logs.Error("Get Deployment failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	if len(DeploymentList) != 0 {
		updateDeployment := DeploymentList[0]
		updateDeployment.Name = this.Name
		updateDeployment.KMetaData = this.KMetaData
		updateDeployment.KSpec = this.KSpec
		updateDeployment.KStatus = this.KStatus
		updateDeployment.NameSpaceName = this.NameSpaceName
		updateDeployment.ClusterName = this.ClusterName
		return updateDeployment.Update()
	} else {
		_, err = o.Insert(this)
		if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
			ResultData.Message = err.Error()
			ResultData.Code = utils.AddDeploymentErr
			logs.Error("Add Deployment failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
			return ResultData
		}
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Deployment) List(from, limit int) Result {
	o := orm.NewOrm()
	var DeploymentList []*Deployment = nil
	var ResultData Result
	var err error
	cond := orm.NewCondition()
	if this.Name != "" {
		cond = cond.And("name__icontains", this.Name)
	}
	if this.HostName != "" {
		cond = cond.And("host_name", this.HostName)
	}
	if this.NameSpaceName != "" {
		cond = cond.And("name_space_name", this.NameSpaceName)
	}
	_, err = o.QueryTable(utils.Deployment).SetCond(cond).Limit(limit, from).All(&DeploymentList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetDeploymentErr
		logs.Error("Get Deployment List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.Deployment).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = DeploymentList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *Deployment) Update() Result {
	o := orm.NewOrm()

	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditDeploymentErr
		logs.Error("Update Deployment: %s failed, code: %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Deployment) Delete() Result {
	o := orm.NewOrm()

	var ResultData Result
	cond := orm.NewCondition()

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if this.ClusterName != "" {
		cond = cond.And("cluster_name", this.ClusterName)
	}
	_, err := o.QueryTable(utils.Deployment).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteDeploymentErr
		logs.Error("Delete Deployment failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}
