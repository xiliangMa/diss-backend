package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type HostConfigInterface interface {
	Add()
	Delete()
	Edit()
	Get()
	List()
	Count()
	GetBnechMarkProportion()
}

func (this *HostConfig) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var HostConfigList []*HostConfig = nil
	var ResultData Result
	var err error
	cond := orm.NewCondition()
	if this.Diss != int8(Diss_All) {
		cond = cond.And("diss", this.Diss)
	}
	if this.DissStatus != int8(Diss_Status_All) {
		cond = cond.And("diss_status", this.DissStatus)
	}
	if this.Label != "" {
		cond = cond.And("Label__contains", this.Label)
	}
	if this.HostName != "" {
		cond = cond.And("host_name__contains", this.HostName)
	}

	if this.AccountName != "" && this.AccountName != Account_Admin {
		cond = cond.And("account_name", this.AccountName)
	}
	_, err = o.QueryTable(utils.HostConfig).SetCond(cond).Limit(limit, from).All(&HostConfigList)
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
