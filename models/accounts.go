package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

// 租户
type Accounts struct {
	Name        string `description:"(名称)"`
	State       string `description:"(状态)"`
	Type        string `description:"(类型)"`
	Email       string `description:"(邮箱)"`
	CreatedAt   int64  `description:"(创建时间)"`
	LastUpdated int64  `description:"(更新时间)"`
}

type AccounstInterface interface {
	List(from, limit int) Result
}

func (this *Accounts) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Diss_Api)
	var accountsList []*Accounts
	var ResultData Result
	var total int64
	var err error
	ignoreAccount := "anchore-system"
	if this.Name == Account_Admin {
		total, err = o.Raw("select * from "+utils.Accounts+" where name != ?", ignoreAccount).QueryRows(&accountsList)
	} else {
		total, err = o.Raw("select * from "+utils.Accounts+" where name = ?", this.Name).QueryRows(&accountsList)
	}

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetAccountsErr
		logs.Error("Get Accounts List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = accountsList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}
