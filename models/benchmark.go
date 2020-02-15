package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type BenchMark struct {
	Id          string `orm:"pk;description(基线id)"`
	Name        string `orm:"description(名称)"`
	Description string `orm:"description(描述)"`
	Type        int8   `orm:"description(类型 docker kubernetes)"`
	path        string `orm:"null;description(模版路径)"`
}

func init() {
	orm.RegisterModel(new(BenchMark))
}

type BenchMarkInterface interface {
	Add()
	Delete()
	Edit()
	Get()
	List()
}

func (this *BenchMark) Add() Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData Result

	_, err := o.Insert(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddBenchMarkErr
		logs.Error("Add BenchMark failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *BenchMark) List() Result {
	o := orm.NewOrm()
	orm.DefaultTimeLoc = time.Local
	o.Using("default")
	var BenchMarkList []*BenchMark
	var ResultData Result

	_, err := o.QueryTable(utils.BenchMark).All(&BenchMarkList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetBenchMarkErr
		logs.Error("Get BenchMark List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.Cluster).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = BenchMarkList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}
