package k8s

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type Image struct {
	Id         string    `orm:"pk;description(镜像id)"`
	Name       string    `orm:"unique;description(镜像名)"`
	Size       int64     `orm:"description(大小)"`
	GroupId    string    `orm:"default(null);description(组id)"`
	GroupName  string    `orm:"default(null);description(组名)"`
	DissStatus int8      `orm:"description(安全状态)"`
	CreateTime time.Time `orm:"description(创建时间);auto_now_add;type(datetime)"`
	UpdateTime time.Time `orm:"null;description(更新时间);auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Image))
}

type ImageInterface interface {
	Add()
	Delete()
	Edit()
	Get()
	List()
}

func (this *Image) Add() models.Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData models.Result

	_, err := o.Insert(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddImageErr
		logs.Error("Add Image failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Image) List() models.Result {
	o := orm.NewOrm()
	orm.DefaultTimeLoc = time.Local
	o.Using("default")
	var ImageList []*Image
	var ResultData models.Result

	_, err := o.QueryTable(utils.Image).All(&ImageList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetImageErr
		logs.Error("Get Image List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.Image).Count()
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
