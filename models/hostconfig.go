package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type HostConfigInterface interface {
	Add()
	Delete()
	Edit()
	Get()
	List()
}

func (this *HostConfig) List(from, limit int) Result {
	o := orm.NewOrm()
	orm.DefaultTimeLoc = time.Local
	o.Using("default")
	var HostConfigList []*HostConfig = nil
	var total = 0
	var ResultData Result
	var err error
	cond := orm.NewCondition()
	if this.HostName != "" {
		cond = cond.And("host_name__contains", this.HostName)
	}

	_, err = o.QueryTable(utils.HostConfig).SetCond(cond).Limit(limit, from).All(&HostConfigList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetHostConfigErr
		logs.Error("GetHostConfig failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	if HostConfigList != nil {
		total = len(HostConfigList)
	}
	data := make(map[string]interface{})
	data["items"] = HostConfigList
	data["total"] = total

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *HostConfig) Update() Result {
	o := orm.NewOrm()
	o.Using("default")
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
