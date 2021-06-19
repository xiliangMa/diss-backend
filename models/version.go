package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type Version struct {
	Id         int64  `orm:"pk;auto" description:"(id)"`
	Name       string `orm:"size(32)" description:"(库名)"`
	Version    string `orm:"size(32)" description:"(版本)"`
	CreateTime int64  `orm:"default(0)" description:"(创建时间)"`
}

func (this *Version) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	this.CreateTime = time.Now().UnixNano()
	_, err := o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddVersionErr
		logs.Error("Add Version failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Version) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var versionList []*Version
	var ResultData Result
	var err error

	total, _ := o.QueryTable(utils.Version).Count()

	if total == 0 {
		v := new(Version)
		v.Name = "漏洞库"
		v.Version = "0.1"
		v.Add()
		v = new(Version)
		v.Name = "病毒库"
		v.Version = "0.1"
		v.Add()
	}

	_, err = o.QueryTable(utils.Version).RelatedSel().Limit(limit, from).OrderBy("-create_time").All(&versionList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetVersionErr
		logs.Error("Get Version List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ = o.QueryTable(utils.Version).Count()
	data := make(map[string]interface{})
	data[Result_Total] = total
	data[Result_Items] = versionList

	ResultData.Code = http.StatusOK
	ResultData.Data = versionList
	return ResultData

}
