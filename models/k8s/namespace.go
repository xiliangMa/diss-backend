package k8s

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type NameSpace struct {
	Id   string `orm:"pk;description(命名空间id)"`
	Name string `orm:"unique;description(命名空间)"`
}

type NameSpaceInterface interface {
	Add()
	Delete()
	Edit()
	Get()
	List()
}

func init() {
	orm.RegisterModel(new(NameSpace))
}

func (this *NameSpace) Add() models.Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData models.Result

	_, err := o.Insert(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddHostErr
		logs.Error("Add NameSpace failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}
