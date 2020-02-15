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
	if this.HostName != "" {
		_, err = o.QueryTable(utils.HostConfig).Filter("host_name", this.HostName).Limit(limit, from).All(&HostConfigList)
	} else {
		_, err = o.QueryTable(utils.HostConfig).Limit(limit, from).All(&HostConfigList)
	}

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetHostListErr
		logs.Error("GetHostList failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
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
