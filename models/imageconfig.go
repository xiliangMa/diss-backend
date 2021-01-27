package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type ImageConfig struct {
	Id            string  `orm:"pk;" description:"(镜像id   k8s拿不到镜像id, 用主机id+镜像名称填充)"`
	ImageId       string  `orm:"" description:"(镜像id)"`
	HostId        string  `orm:"" description:"(主机id)"`
	HostName      string  `orm:"" description:"(主机名称)"`
	Name          string  `orm:"" description:"(镜像名)"`
	Size          string  `orm:"" description:"(大小)"`
	OS            string  `orm:"" description:"(镜像名)"`
	DissStatus    int8    `orm:"" description:"(安全状态)"`
	Age           string  `orm:"default(null);" description:"(运行时长)"`
	CreateTime    int64   `orm:"default(0)" description:"(创建时间)"`
	DBType        string  `-" description:"(数据库类型 Mysql Oracle Redis Postgres Mongodb Memcache DB2 Hbase)"`
	GetLatestTask bool    `-" description:"(是否获取最新一个task、否则获取所有task列表)"`
	TaskList      []*Task `orm:"reverse(many);null" description:"(任务列表)"`
}

type ImageConfigInterface interface {
	Add() Result
	Delete() Result
	Get() *ImageConfig
	List(from, limit int) Result
	GetDBImageByType() Result
}

func (this *ImageConfig) Get() *ImageConfig {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	object := new(ImageConfig)
	var err error
	cond := orm.NewCondition()
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if this.ImageId != "" {
		cond = cond.And("image_id", this.ImageId)
	}

	err = o.QueryTable(utils.ImageConfig).SetCond(cond).RelatedSel().One(object)
	if err != nil {
		logs.Error("Get ImageConfig failed, code: %d, err: %s", err.Error(), utils.GetImageContentErr)
		return nil
	}
	return object
}

func (this *ImageConfig) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	var err error
	var imageConfigList []*ImageConfig

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
	_, err = o.QueryTable(utils.ImageConfig).SetCond(cond).All(&imageConfigList)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetContainerConfigErr
		logs.Error("Get ContainerConfig failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	if len(imageConfigList) != 0 {
		// agent 或者 k8s 数据更新（因为没有diss-backend的关系数据，所以直接更新）
		return this.Update()
	} else {
		_, err = o.Insert(this)
		if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
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
	o.Using(utils.DS_Default)
	var imageConfigList []*ImageConfig
	var ResultData Result
	var err error

	cond := orm.NewCondition()
	if this.HostId != "" {
		cond = cond.And("host_id", this.HostId)
	}
	if this.HostName != "" {
		cond = cond.And("host_name__icontains", this.HostName)
	}
	if this.ImageId != "" {
		cond = cond.And("image_id__contains", this.ImageId)
	}
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if this.Name != "" {
		cond = cond.And("name__contains", this.Name)
	}

	_, err = o.QueryTable(utils.ImageConfig).RelatedSel().SetCond(cond).Limit(limit, from).All(&imageConfigList)
	for _, image := range imageConfigList {
		if this.GetLatestTask {
			o.LoadRelated(image, "TaskList", 1, 1, 0, "-update_time")
		} else {
			o.LoadRelated(image, "TaskList", 1, limit, 0, "-update_time")
		}
	}

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetImageConfigErr
		logs.Error("Get ImageConfig List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.ImageConfig).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = imageConfigList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *ImageConfig) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
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

func (this *ImageConfig) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()

	// 根据agent同步时 依据 host_id 删除该主机上所有的容器历史记录
	if this.HostId != "" {
		cond = cond.And("host_id", this.HostId)
	}
	_, err := o.QueryTable(utils.ImageConfig).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteImageInfoErr
		logs.Error("Delete ImageConfig failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}

func (this *ImageConfig) GetDBCountByType() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	ResultData := Result{Code: http.StatusOK}
	var dbCount []orm.Params
	var err error
	var total int64 = 0

	cond := orm.NewCondition()
	if this.HostId != "" {
		cond = cond.And("host_id", this.HostId)
	}

	_, err = o.Raw(utils.GetDBCountSql(this.HostId)).Values(&dbCount)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetImageConfigErr
		logs.Error("Get ImageConfig count group type failed, code: %d, err: %s", utils.GetImageConfigErr, err.Error())
		return ResultData
	}

	total, _ = o.QueryTable(utils.ImageConfig).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = dbCount

	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *ImageConfig) GetDBImageByType() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var imageList []*ImageConfig
	var ResultData Result
	var err error
	cond := orm.NewCondition()
	if this.HostId != "" {
		cond = cond.And("host_id", this.HostId)
	}
	if this.DBType != "" {
		cond = cond.And("name__icontains", this.DBType)
	}
	_, err = o.QueryTable(utils.ImageConfig).SetCond(cond).All(&imageList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetImageConfigErr
		logs.Error("Get ImageConfig List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	total, _ := o.QueryTable(utils.ImageConfig).SetCond(cond).Count()
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
