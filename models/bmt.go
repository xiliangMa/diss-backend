package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type BenchMarkTemplate struct {
	Id          string `orm:"pk;" description:"(基线id)"`
	Name        string `orm:"" description:"(名称)"`
	Description string `orm:"" description:"(描述)"`
	Type        int8   `orm:"" description:"(类型 docker 0  kubernetes 1)"`
	Path        string `orm:"null;" description:"(模版路径)"`
	Commands    string `orm:"null;" description:"(操作命令)"`
	Status      int8   `orm:"default(0);" description:"(类型 停用 0  启用 1)"`
}

type BenchMarkInterface interface {
	Add()
	Delete()
	Edit()
	Get()
	List()
}

func (this *BenchMarkTemplate) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddBenchMarkTemplateErr
		logs.Error("Add BenchMarkTemplate failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *BenchMarkTemplate) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var BenchMarkTemplateList []*BenchMarkTemplate
	var ResultData Result
	var err error

	// to do Multiple conditions search
	if this.Name == "" {
		_, err = o.QueryTable(utils.BenchMarkTemplate).Limit(limit, from).All(&BenchMarkTemplateList)
	} else {
		_, err = o.QueryTable(utils.BenchMarkTemplate).Filter("name__icontains", this.Name).Limit(limit, from).All(&BenchMarkTemplateList)
	}
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetBenchMarkTemplateErr
		logs.Error("Get BenchMarkTemplate List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.BenchMarkTemplate).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = BenchMarkTemplateList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}
