package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type ImageBlocking struct {
	Id          string `orm:"pk;size(128)" description:"(Id)"`
	ImageId     string `orm:"" description:"(镜像id)"`
	ImageName   string `orm:"" description:"(镜像名)"`
	AccountName string `orm:"" description:"(租户)"`
	CreateTime  int64  `orm:"default(0)" description:"(阻断时间)"`
	Status      string `orm:"default(null);" description:"(状态)"`
	Action      string `orm:"-" description:"(处理方式：kill)"`
}

type ImageBlockingInterface interface {
	Add() Result
	Get() *Result
	GetIB() *ImageBlocking
	List(from, limit int) Result
}

func (this *ImageBlocking) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	var err error

	if this.Id == "" {
		uid, _ := uuid.NewV4()
		this.Id = uid.String()
	}

	if this.CreateTime == 0 {
		this.CreateTime = time.Now().UnixNano()
	}

	_, err = o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddImageBlockingErr
		logs.Error("Add ImageBlocking failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *ImageBlocking) Get() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()

	if this.ImageId != "" {
		cond = cond.And("image_id", this.ImageId)
	}

	if this.ImageName != "" {
		cond = cond.And("image_name", this.ImageName)
	}

	total, _ := o.QueryTable(utils.ImageBlocking).SetCond(cond).RelatedSel().Count()

	data := make(map[string]interface{})
	data[Result_Total] = total

	ResultData.Code = http.StatusOK
	ResultData.Data = data

	return ResultData
}

func (this *ImageBlocking) GetIb() *ImageBlocking {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	imageBlocking := new(ImageBlocking)
	cond := orm.NewCondition()

	if this.ImageId != "" {
		cond = cond.And("image_id", this.ImageId)
	}

	if this.ImageName != "" {
		cond = cond.And("image_name", this.ImageName)
	}

	err := o.QueryTable(utils.ImageBlocking).SetCond(cond).RelatedSel().One(imageBlocking)

	if err != nil {
		return nil
	}

	return imageBlocking
}

func (this *ImageBlocking) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var imageBlockingList []*ImageBlocking
	var ResultData Result
	var err error

	cond := orm.NewCondition()

	if this.AccountName != "" && this.AccountName != Account_Admin {
		cond = cond.And("account_name", this.AccountName)
	}

	if this.ImageName != "" {
		cond = cond.And("image_name__contains", this.ImageName)
	}

	_, err = o.QueryTable(utils.ImageBlocking).RelatedSel().SetCond(cond).Limit(limit, from).All(&imageBlockingList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetImageBlockingErr
		logs.Error("Get ImageBlocking List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
	}

	total, _ := o.QueryTable(utils.ImageBlocking).SetCond(cond).Count()
	data := make(map[string]interface{})
	data[Result_Total] = total
	data[Result_Items] = imageBlockingList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}
