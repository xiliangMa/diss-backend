package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type ImageConfig struct {
	Id         string    `orm:"pk;description(镜像id   k8s拿不到镜像id, 用主机id+镜像名称填充)"`
	ImageId    string    `orm:"description(镜像id)"`
	HostId     string    `orm:"description(主机id)"`
	Name       string    `orm:"description(镜像名)"`
	Size       string    `orm:"description(大小)"`
	OS         string    `orm:"description(镜像名)"`
	DissStatus int8      `orm:"description(安全状态)"`
	Age        string    `orm:"default(null);description(运行时长)"`
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
	var err error
	var imageConfiggList []*ImageConfig

	cond := orm.NewCondition()
	if this.HostId != "" {
		cond = cond.And("host_id", this.HostId)
	}
	if this.ImageId != "" {
		cond = cond.And("image_id", this.ImageId)
	}
	if this.Name != "" {
		cond = cond.And("name", this.Name)
	}
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	_, err = o.QueryTable(utils.ImageConfig).SetCond(cond).All(&imageConfiggList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetContainerConfigErr
		logs.Error("Get ContainerConfig failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	if len(imageConfiggList) != 0 {
		// agent 或者 k8s 数据更新（因为没有diss-backend的关系数据，所以直接更新）
		return this.Update()
	} else {
		_, err = o.Insert(this)
		if err != nil {
			ResultData.Message = err.Error()
			ResultData.Code = utils.AddImageConfigErr
			logs.Error("Add ImageConfig failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
			return ResultData
		}
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *ImageConfig) List(from, limit int) Result {
	o := orm.NewOrm()
	orm.DefaultTimeLoc = time.Local
	o.Using("default")
	var imageList []*ImageConfig
	var ResultData Result
	var err error
	var total = 0

	cond := orm.NewCondition()
	if this.HostId != "" {
		cond = cond.And("host_id", this.HostId)
	}
	if this.ImageId != "" {
		cond = cond.And("image_id", this.ImageId)
	}
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if this.Name != "" {
		cond = cond.And("name__contains", this.Name)
	}

	_, err = o.QueryTable(utils.ImageConfig).SetCond(cond).Limit(limit, from).All(&imageList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetImageConfigErr
		logs.Error("Get ImageConfig List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	if imageList != nil {
		total = len(imageList)
	}
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = imageList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *ImageConfig) Update() Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditImageConfigErr
		logs.Error("Update ImageConfig: %s failed, code: %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}
