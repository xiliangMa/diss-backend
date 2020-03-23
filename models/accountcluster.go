package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type AccountCluster struct {
	Id          string `orm:"pk" description:"(uuid)`
	AccountName string `description:"(租户)"`
	ClusterId   string `description:"(集群id)"`
}

func init() {
	orm.RegisterModel(new(AccountCluster))
}

type AccountClusterInterface interface {
	List()
	Add()
}

func (this *AccountCluster) List() (error, []string) {
	o := orm.NewOrm()
	orm.DefaultTimeLoc = time.Local
	o.Using(utils.DS_Default)
	var cIds []string

	_, err := o.Raw("select cluster_id from "+utils.AccountCluster+" where account_name = ?", this.AccountName).QueryRows(&cIds)
	if err != nil {
		logs.Error("Get AccountCluster List failed, code: %d, err: %s", utils.GetAccountClusterErr, err)
	}
	return err, cIds
}

func (this *AccountCluster) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddAccountClusterErr
		logs.Error("Add AccountCluster failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}
