package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type ImageInfo struct {
	Id          string `orm:"pk;description(数据库主键)"`
	ImageId     string `orm:"description(镜像id)"`
	Name        string `orm:"description(镜像名)"`
	HostId      string `orm:"description(主机id)"`
	RepoTags    string `orm:"description(RepoTags)"`
	RepoDigests string `orm:"description(RepoDigests)"`
	Os          string `orm:"description(系统)"`
	Created     string `orm:"description(创建时间)"`
	Size        string `orm:"description(大小)"`
	Layers      string `orm:"description(Layers)"`
}

func init() {
	orm.RegisterModel(new(ImageInfo))
}

type ImageInfoInterface interface {
	Add()
	Delete()
	Edit()
	Get()
	List()
}

func (this *ImageInfo) Add() Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData Result
	var err error
	var imageInfoList []*ImageInfo

	cond := orm.NewCondition()

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}

	if this.HostId != "" {
		cond = cond.And("host_id", this.HostId)
	}
	if this.ImageId != "" {
		cond = cond.And("image_id", this.ImageId)
	}

	_, err = o.QueryTable(utils.ImageInfo).SetCond(cond).All(&imageInfoList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetImageInfoErr
		logs.Error("Get ImageInfo failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	if len(imageInfoList) != 0 {
		// agent 或者 k8s 数据更新（因为没有diss-backend的关系数据，所以直接删除在添加）
		if result := this.Delete(); result.Code != http.StatusOK {
			return result
		}
	}
	_, err = o.Insert(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddImageInfoErr
		logs.Error("Add ImageInfo failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *ImageInfo) List() Result {
	o := orm.NewOrm()
	orm.DefaultTimeLoc = time.Local
	o.Using("default")
	var imageList []*ImageInfo
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

	_, err = o.QueryTable(utils.ImageInfo).SetCond(cond).All(&imageList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetImageInfoErr
		logs.Error("Get ImageInfo List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
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

func (this *ImageInfo) Delete() Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData Result
	_, err := o.Delete(&ImageInfo{Id: this.Id})
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteImageInfoErr
		logs.Error("Delete ImageInfo failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}
