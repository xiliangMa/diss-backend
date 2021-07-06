package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ImageConfig struct {
	Id         string    `orm:"pk;" description:"(镜像id   k8s拿不到镜像id, 用主机id+镜像名称填充)"`
	ImageId    string    `orm:"" description:"(镜像id)"`
	HostId     string    `orm:"" description:"(主机id)"`
	HostName   string    `orm:"" description:"(主机名称)"`
	Name       string    `orm:"" description:"(镜像名)"`
	Size       string    `orm:"" description:"(大小)"`
	OS         string    `orm:"column(os)" description:"(镜像名)"`
	DissStatus int8      `orm:"" description:"(安全状态)"`
	Age        string    `orm:"default(null);" description:"(运行时长)"`
	CreateTime int64     `orm:"default(0)" description:"(创建时间)"`
	TaskList   []*Task   `orm:"reverse(many);null" description:"(任务列表)"`
	Registry   *Registry `orm:"rel(fk);default(null);null" description:"(仓库)"`
	Type       string    `orm:"-" description:"(区分主机镜像还是仓库镜像)"`
	Namespaces string    `orm:"-" description:"(命名空间)"`
	VirusScan  bool      `orm:"-" description:"(区分杀毒还是扫描)"`
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
	if this.HostId != "" {
		cond = cond.And("host_id", this.HostId)
	}

	if this.ImageId != "" {
		cond = cond.And("image_id", this.ImageId)
	}

	if this.Name != "" {
		cond = cond.And("name", this.Name)
	}

	if this.Registry != nil {
		cond = cond.And("registry_id", this.Registry.Id)
	}

	err = o.QueryTable(utils.ImageConfig).SetCond(cond).RelatedSel().One(object)
	if err != nil {
		return nil
	}
	return object
}

func (this *ImageConfig) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	var err error

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
	if this.Registry != nil {
		cond = cond.And("registry_id", this.Registry.Id)
	}
	if this.Id == "" {
		this.Id = strconv.FormatInt(this.Registry.Id, 10) + "---" + this.ImageId + "---" + this.Name
	}
	if this.CreateTime == 0 {
		this.CreateTime = time.Now().UnixNano()
	}

	_, err = o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddImageConfigErr
		logs.Error("Add ImageConfig failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
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
	// 根据type类型判断是主机镜像（type不为空）还是仓库镜像
	if this.Type != "" {
		cond = cond.And("registry_id__isnull", true)
	} else {
		cond = cond.And("registry_id__isnull", false)
	}
	if this.Registry != nil {
		cond = cond.And("registry_id", this.Registry.Id)
	}

	_, err = o.QueryTable(utils.ImageConfig).RelatedSel().SetCond(cond).OrderBy("-create_time").Limit(limit, from).All(&imageConfigList)
	for _, image := range imageConfigList {
		cond2 := orm.NewCondition()
		var task []*Task
		if this.VirusScan {
			cond2 = cond2.And("image_id", image.Id)
			cond2 = cond2.AndNot("virus_status", "")
			_, err = o.QueryTable(utils.Task).RelatedSel().SetCond(cond2).OrderBy("-update_time").Limit(1, 0).All(&task)
			image.TaskList = task
		} else {
			cond2 = cond2.And("image_id", image.Id)
			cond2 = cond2.AndNot("scan_status", "")
			_, err = o.QueryTable(utils.Task).RelatedSel().SetCond(cond2).OrderBy("-update_time").Limit(1, 0).All(&task)
			image.TaskList = task
		}
	}

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetImageConfigErr
		logs.Error("Get ImageConfig List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
	}

	total, _ := o.QueryTable(utils.ImageConfig).SetCond(cond).Count()
	data := make(map[string]interface{})
	data[Result_Total] = total
	data[Result_Items] = imageConfigList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
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
	if this.Id != "" {
		cond = cond.And("id__in", strings.Split(this.Id, ","))
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

	var fields []string
	if this.HostId != "" {
		fields = append(fields, this.HostId, this.HostId, this.HostId, this.HostId, this.HostId, this.HostId, this.HostId, this.HostId)
	}
	_, err = o.Raw(utils.GetDBCountSql(this.HostId), fields).Values(&dbCount)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetImageConfigErr
		logs.Error("Get ImageConfig count group type failed, code: %d, err: %s", utils.GetImageConfigErr, err.Error())
		return ResultData
	}

	total, _ = o.QueryTable(utils.ImageConfig).SetCond(cond).Count()
	data := make(map[string]interface{})
	data[Result_Total] = total
	data[Result_Items] = dbCount

	ResultData.Data = data
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
	_, err = o.QueryTable(utils.ImageConfig).SetCond(cond).All(&imageList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetImageConfigErr
		logs.Error("Get ImageConfig List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	total, _ := o.QueryTable(utils.ImageConfig).SetCond(cond).Count()
	data := make(map[string]interface{})
	data[Result_Total] = total
	data[Result_Items] = imageList
	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *ImageConfig) Count() int64 {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	cond := orm.NewCondition()
	if this.Type != All {
		if this.Type != "" {
			cond = cond.And("registry_id__isnull", true)
		} else {
			cond = cond.And("registry_id__isnull", false)
		}
	}
	count, _ := o.QueryTable(utils.ImageConfig).SetCond(cond).Count()
	return count

}
