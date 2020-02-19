package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type HostInfoInterface interface {
	Add()
	Delete()
	Edit()
	Get()
	List()
}

func (this *HostInfo) List(id string, from, limit int) Result {
	o := orm.NewOrm()
	orm.DefaultTimeLoc = time.Local
	o.Using("default")
	var HostInfoList []*HostInfo = nil
	var ResultData Result
	var total = 0

	_, err := o.QueryTable(utils.HostInfo).Filter("id", id).Limit(limit, from).All(&HostInfoList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetHostInfoErr
		logs.Error("GetHostInfoList failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	if HostInfoList != nil {
		total = len(HostInfoList)
	}
	data := make(map[string]interface{})
	data["items"] = HostInfoList
	data["total"] = total

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *HostInfo) Update() Result {
	o := orm.NewOrm()
	o.Using("default")
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
