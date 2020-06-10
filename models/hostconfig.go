package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type HostConfigInterface interface {
	Inner_AddHostConfig() error
	Inner_AddHostInfo() error
	List(from, limit int) Result
	Update() Result
	Delete() Result
	UpdateDynamic() Result
	Count() int64
	GetBnechMarkProportion() (int64, int64)
}

func (this *HostConfig) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var HostConfigList []*HostConfig = nil
	var ResultData Result
	var err error
	cond := orm.NewCondition()
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if this.Diss != "" && this.Diss != All {
		cond = cond.And("diss", this.Diss)
	}
	if this.DissStatus != "" && this.DissStatus != All {
		cond = cond.And("diss_status", this.DissStatus)
	}
	if this.Label != "" {
		cond = cond.And("label__contains", this.Label)
	}
	if this.GroupId != "" {
		cond = cond.And("Group", this.GroupId)
	}
	if this.HostName != "" {
		cond = cond.And("host_name__contains", this.HostName)
	}
	if this.AccountName != "" && this.AccountName != Account_Admin {
		cond = cond.And("account_name", this.AccountName)
	}
	_, err = o.QueryTable(utils.HostConfig).SetCond(cond).Limit(limit, from).OrderBy("-host_name").RelatedSel().All(&HostConfigList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetHostConfigErr
		logs.Error("GetHostConfig failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.HostConfig).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["items"] = HostConfigList
	data["total"] = total
	if total == 0 {
		ResultData.Data = nil
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *HostConfig) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditHostConfigErr
		logs.Error("Update HostConfig: %s failed, code: %d, err: %s", this.HostName, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *HostConfig) UpdateDynamic() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	hostConfig := new(HostConfig)
	cond := orm.NewCondition()
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if err := o.QueryTable(utils.HostConfig).SetCond(cond).One(hostConfig); err != nil {
		ResultData.Code = utils.HostConfigNotFoundErr
		ResultData.Message = err.Error()
		logs.Error("Get HostConfig: %s failed, code: %d, err: %s", this.HostName, ResultData.Code, ResultData.Message)
	}
	hostConfig.PublicAddr = this.PublicAddr
	_, err := o.Update(hostConfig)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditHostConfigDynamicErr
		logs.Error("Update HostInfo Dynamic failed, HostName: %s, failed, code: %d, err: %s", this.HostName, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *HostConfig) Count() int64 {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	count, _ := o.QueryTable(utils.HostConfig).Count()
	return count
}

// docker基线 / k8s 基线
func (this *HostConfig) GetBnechMarkProportion() (int64, int64) {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	dockerBenchMarkCount, _ := o.QueryTable(utils.HostConfig).Count()
	k8sBenchMarkCount, _ := o.QueryTable(utils.HostConfig).Filter("is_in_k8s", false).Count()
	return dockerBenchMarkCount, k8sBenchMarkCount
}

func (this *HostConfig) Delete() Result {
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
	_, err := o.QueryTable(utils.HostConfig).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteHostConfigErr
		logs.Error("Delete HostConfig failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}
