package k8s

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type Container struct {
	Id            string    `orm:"pk;description(id)"`
	Name          string    `orm:"unique;description(容器名)"`
	NameSpaceId   string    `orm:"description(命名空间id)"`
	NameSpaceName string    `orm:"description(命名空间)"`
	Status        uint8     `orm:"default(0);description(状态)"`
	Command       string    `orm:"default(null);description(命令)"`
	Image         string    `orm:"default(null);description(镜像)"`
	CreateTime    time.Time `orm:"description(创建时间);auto_now_add;type(datetime)"`
	UpdateTime    time.Time `orm:"null;description(更新时间);type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Container))
}

type ContainerInterface interface {
	Add()
	Delete()
	Edit()
	Get()
	List()
}


func (this *Container) Add() models.Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData models.Result

	_, err := o.Insert(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddContainerErr
		logs.Error("Add Container failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Container) List() models.Result {
	o := orm.NewOrm()
	orm.DefaultTimeLoc = time.Local
	o.Using("default")
	var ContainerList []*Container
	var ResultData models.Result

	_, err := o.QueryTable(utils.Container).All(&ContainerList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetContainerErr
		logs.Error("Get Container List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.Container).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = ContainerList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}
