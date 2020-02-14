package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type ImageConfig struct {
	Id         string    `orm:"pk;description(镜像id)"`
	HostId     string    `orm:"description(主机id)"`
	Name       string    `orm:"unique;description(镜像名)"`
	Size       int64     `orm:"description(大小)"`
	OS         string    `orm:"description(镜像名)"`
	DissStatus int8      `orm:"description(安全状态)"`
	CreateTime time.Time `orm:"null;description(创建时间);type(datetime)"`
}

func init() {
	orm.RegisterModel(new(ImageConfig))
}

type ImageConfigInterface interface {
	Add()
	Delete()
	Edit()
	Get()
	List()
}

func (this *ImageConfig) Add() Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData Result

	_, err := o.Insert(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddImageConfigErr
		logs.Error("Add ImageConfig failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *ImageConfig) List() Result {
	o := orm.NewOrm()
	orm.DefaultTimeLoc = time.Local
	o.Using("default")
	var ImageList []*ImageConfig
	var ResultData Result

	_, err := o.QueryTable(utils.ImageConfig).All(&ImageList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetImageConfigErr
		logs.Error("Get ImageConfig List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.ImageConfig).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = ImageList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}
