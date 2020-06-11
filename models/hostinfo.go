package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type HostInfoInterface interface {
	Inner_AddHostInfo() error
	Update() Result
	List() Result
	Delete() Result
	UpdateDynamic() Result
}

func (this *HostInfo) List() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var HostInfoList []*HostInfo = nil
	var ResultData Result
	cond := orm.NewCondition()
	if this.HostName != "" {
		cond = cond.And("host_name__contains", this.HostName)
	}
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}

	_, err := o.QueryTable(utils.HostInfo).SetCond(cond).All(&HostInfoList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetHostInfoErr
		logs.Error("GetHostInfoList failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.HostInfo).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["items"] = HostInfoList
	data["total"] = total

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *HostInfo) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditHostInfoErr
		logs.Error("Update HostInfo: %s failed, code: %d, err: %s", this.HostName, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *HostInfo) UpdateDynamic() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	hostInfo := new(HostInfo)
	cond := orm.NewCondition()
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if err := o.QueryTable(utils.HostInfo).SetCond(cond).One(hostInfo); err != nil || hostInfo != nil {
		ResultData.Code = utils.HostConfigNotFoundErr
		ResultData.Message = err.Error()
		logs.Error("Get HostInfo: %s failed, code: %d, err: %s", this.HostName, ResultData.Code, ResultData.Message)
	}
	hostInfo.PublicAddr = this.PublicAddr
	hostInfo.ImageCount = this.ImageCount
	hostInfo.ContainerCount = this.ContainerCount
	hostInfo.ContainerRunningCount = this.ContainerRunningCount
	hostInfo.ContainerPausedCount = this.ContainerPausedCount
	hostInfo.ContainerStoppedCount = this.ContainerStoppedCount
	_, err := o.Update(hostInfo)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditHostInfoDynamicErr
		logs.Error("Update HostInfo Dynamic failed, HostName: %s, code: %d, err: %s", this.HostName, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *HostInfo) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if this.ClusterId != "" {
		cond = cond.And("cluster_id", this.ClusterId)
	}
	_, err := o.QueryTable(utils.HostInfo).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteHostinfoErr
		logs.Error("Delete HostInfo failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}
