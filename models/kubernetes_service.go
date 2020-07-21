package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type Service struct {
	Id            string `orm:"pk;" description:"(Service id)"`
	Name          string `orm:"size(128)" description:"(Service名)"`
	NameSpaceName string `orm:"size(255);default(null);" description:"(命名空间)"`
	ClusterId     string `orm:"size(128)" description:"(集群Id)"`
	ClusterName   string `orm:"size(32)" description:"(集群名)"`
	KMetaData     string `orm:"" description:"(源数据)"`
	KSpec         string `orm:"" description:"(Spec数据)"`
	KStatus       string `orm:"" description:"(状态数据)"`
}

type ServiceInterface interface {
	Add() Result
	Delete() Result
	Update() Result
	List(from, limit int) Result
	EmptyDirtyData() error
}

func (this *Service) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	var ServiceList []*Service
	var err error
	cond := orm.NewCondition()
	if this.Name != "" {
		cond = cond.And("id", this.Id)
	}

	_, err = o.QueryTable(utils.Service).SetCond(cond).All(&ServiceList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetServiceErr
		logs.Error("Get Service failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	if len(ServiceList) != 0 {
		updateService := ServiceList[0]
		updateService.Name = this.Name
		updateService.KMetaData = this.KMetaData
		updateService.KSpec = this.KSpec
		updateService.KStatus = this.KStatus
		updateService.NameSpaceName = this.NameSpaceName
		updateService.ClusterName = this.ClusterName
		return updateService.Update()
	} else {
		_, err = o.Insert(this)
		if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
			ResultData.Message = err.Error()
			ResultData.Code = utils.AddServiceErr
			logs.Error("Add Service failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
			return ResultData
		}
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Service) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ServiceList []*Service = nil
	var ResultData Result
	var err error
	cond := orm.NewCondition()
	if this.Id != "" {
		cond = cond.And("id", this.Name)
	}
	if this.Name != "" {
		cond = cond.And("name__icontains", this.Name)
	}
	if this.NameSpaceName != "" {
		cond = cond.And("name_space_name", this.NameSpaceName)
	}
	_, err = o.QueryTable(utils.Service).SetCond(cond).Limit(limit, from).All(&ServiceList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetServiceErr
		logs.Error("Get Service List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.Service).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = ServiceList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *Service) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditServiceErr
		logs.Error("Update Service: %s failed, code: %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Service) Delete() Result {
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
	if this.ClusterId != "" {
		cond = cond.And("cluster_id", this.ClusterId)
	}
	_, err := o.QueryTable(utils.Service).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteServiceErr
		logs.Error("Delete Service failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}
