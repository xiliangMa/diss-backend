package models

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type SystemTemplate struct {
	Id                       string                 `orm:"pk;" description:"(id)"`
	Account                  string                 `orm:"default(admin)" description:"(租户)"`
	Name                     string                 `orm:"" description:"(名称)"`
	Description              string                 `orm:"" description:"(描述)"`
	Type                     string                 `orm:"" description:"(类型)"`
	Version                  string                 `orm:"null" description:"(版本)"`
	Commands                 string                 `orm:"null;" description:"(操作命令)"`
	Status                   string                 `orm:"default(Enable);" description:"(类型 Enable Disable)"`
	IsDefault                bool                   `orm:"default(false);" description:"(系统策略默认项)"`
	IsSystem                 bool                   `orm:"default(false);" description:"(是否是系统策略)"`
	SystemTemplateGroup      []*SystemTemplateGroup `orm:"rel(m2m);" description:"(策略组)"`
	Job                      []*Job                 `orm:"reverse(many);null" description:"(job)"`
	Task                     []*Task                `orm:"reverse(many);null" description:"(task)"`
	ConfigMode               string                 `orm:"default(all);null" description:"(配置形式，如target、checks、group)"`
	DefaultTargets           string                 `orm:"null;" description:"(默认的target枚举)"`
	CheckMasterJson          string                 `orm:"null;" description:"(Master target的json内容)"`
	CheckNodeJson            string                 `orm:"null;" description:"(Node target的json内容)"`
	CheckControlPlaneJson    string                 `orm:"null;" description:"(ControlPlane target的json内容)"`
	CheckEtcdJson            string                 `orm:"null;" description:"(Etcd target的json内容)"`
	CheckPoliciesJson        string                 `orm:"null;" description:"(Polices target的json内容)"`
	CheckManagedServicesJson string                 `orm:"null;" description:"(ManagedServices target的json内容)"`
	CheckIdsMaster           string                 `orm:"null;" description:"(选中的检查项Id列表-Master)"`
	CheckIdsNode             string                 `orm:"null;" description:"(选中的检查项Id列表-Node)"`
	CheckIdsControlPlane     string                 `orm:"null;" description:"(选中的检查项Id列表-ControlPlane)"`
	CheckIdsEtcd             string                 `orm:"null;" description:"(选中的检查项Id列表-Etcd)"`
	CheckIdsPolicies         string                 `orm:"null;" description:"(选中的检查项Id列表-Policies)"`
	CheckIdsManagedServices  string                 `orm:"null;" description:"(选中的检查项Id列表-ManagedServices)"`
	CheckIdsDocker           string                 `orm:"null;" description:"(选中的检查项Id列表-Docker)"`
	CheckIdsDockerCheck      string                 `orm:"null;" description:"(选中的检查项Id列表-DockerCheck)"`
	CheckIds                 string                 `orm:"null;" description:"(选中的检查项集合)"`
	DefaultPathList          string                 `orm:"null;" description:"(默认的路径集合)"`
	EngineCode               string                 `orm:"null;" description:"(使用的引擎编码)"`
	CreateTime               int64                  `orm:"default(0);" description:"(创建时间)"`
	UpdateTime               int64                  `orm:"default(0)" description:"(更新时间)"`
}

type SystemTemplateGroup struct {
	Id             string            `orm:"pk;" description:"(基线id)"`
	Account        string            `orm:"default(admin)" description:"(租户)"`
	Name           string            `orm:"" description:"(名称)"`
	Description    string            `orm:"" description:"(描述)"`
	Type           string            `orm:"" description:"(类型)"`
	Version        string            `orm:"null" description:"(版本)"`
	Status         string            `orm:"default(Enable);" description:"(类型 Enable Disable)"`
	IsDefault      bool              `orm:"default(false);" description:"(默认系统策略)"`
	SystemTemplate []*SystemTemplate `orm:"reverse(many);null" description:"(策略模板)"`
	Job            []*Job            `orm:"reverse(many);null" description:"(job)"`
}

type SystemTemplateInterface interface {
	Add() Result
	List() Result
	Delete() Result
	Update() Result
	GetDefaultTemplate() map[string]*SystemTemplate
}
type SystemTemplateGroupInterface interface {
	Add() Result
	List() Result
	Delete() Result
	Update() Result
	Get() (*SystemTemplate, error)
}

func (this *SystemTemplate) Get() (*SystemTemplate, error) {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	cond := orm.NewCondition()
	var template SystemTemplate

	if this.Type != "" {
		cond = cond.And("type", this.Type)
	}
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if this.Version != "" {
		cond = cond.And("version", this.Version)
	}
	if this.Name != "" {
		cond = cond.And("name", this.Name)
	}

	err := o.QueryTable(utils.SYSTemplate).SetCond(cond).One(&template)
	if err != nil {
		errMsg := fmt.Sprintf("Get SystemTemplate failed, code: %d, err: %s", utils.GetSYSTemplateErr, err.Error())
		logs.Error(errMsg)
		return &template, err
	}

	return &template, err
}

func (this *SystemTemplate) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	if this.Name == "" || this.Type == "" || this.Version == "" {
		ResultData.Code = utils.AddSYSTemplateErr
		errMsg := fmt.Sprint("Name, Type or Version Not Input for SystemTemplate, code: %s", utils.AddSYSTemplateErr)
		ResultData.Message = errMsg

		logs.Error(errMsg)
		return ResultData
	}

	uuidCode, _ := uuid.NewV4()
	this.Id = uuidCode.String()

	if this.Account == "" {
		this.Account = Account_Admin
	}

	this.UpdateTime = time.Now().UnixNano()
	this.CreateTime = time.Now().UnixNano()

	_, err := o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {

		logs.Error("Add SystemTemplate failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *SystemTemplate) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var systemTemplateList []*SystemTemplate
	var ResultData Result
	var err error
	cond := orm.NewCondition()

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if this.Account != "" {
		cond = cond.And("account", this.Account)
	}
	if this.Name != "" {
		cond = cond.And("name__contains", this.Name)
	}
	if this.Type != "" {
		cond = cond.And("type", this.Type)
	}
	if this.Version != All {
		cond = cond.And("name__contains", this.Name)
	}
	if this.Commands != "" {
		cond = cond.And("commands__contains", this.Commands)
	}
	if this.Status != "" && this.Status != All {
		cond = cond.And("status", this.Status)
	}
	_, err = o.QueryTable(utils.SYSTemplate).SetCond(cond).RelatedSel().Limit(limit, from).OrderBy("-create_time").All(&systemTemplateList)
	for _, systemTemplate := range systemTemplateList {
		o.LoadRelated(systemTemplate, "SystemTemplateGroup")
		o.LoadRelated(systemTemplate, "Job")
		o.LoadRelated(systemTemplate, "Task")
	}
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetSYSTemplateErr
		logs.Error("Get SystemTemplate List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.SYSTemplate).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = systemTemplateList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *SystemTemplate) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	_, err := o.QueryTable(utils.SYSTemplate).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteSYSTemplateErr
		logs.Error("Delete SYSTemplateErr failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}

func (this *SystemTemplate) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	sysTemplateQuery := SystemTemplate{}
	sysTemplateQuery.Id = this.Id
	sysTemplateObj, _ := sysTemplateQuery.Get()

	this.CreateTime = sysTemplateObj.CreateTime
	this.UpdateTime = time.Now().UnixNano()
	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditSYSTemplateErr
		logs.Error("Update SYSTemplateErr: %s failed, code: %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *SystemTemplate) GetDefaultTemplate() map[string]*SystemTemplate {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var systemTemplateList []*SystemTemplate
	defaultTemplateList := make(map[string]*SystemTemplate)
	var err error
	cond := orm.NewCondition()
	cond = cond.And("is_default", true)
	_, err = o.QueryTable(utils.SYSTemplate).SetCond(cond).RelatedSel().All(&systemTemplateList)
	if err != nil {
		logs.Error("Get SystemTemplate List failed, code: %d, err: %s", utils.GetSYSTemplateErr, err.Error())
		return nil
	}

	for _, systemTemplate := range systemTemplateList {
		defaultTemplateList[systemTemplate.Type] = systemTemplate
	}
	return defaultTemplateList
}

func (this *SystemTemplateGroup) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	this.IsDefault = false

	if this.Id == "" {
		uid, _ := uuid.NewV4()
		this.Id = uid.String()
	}
	_, err := o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddSYSTemplateGroupErr
		logs.Error("Add SystemTemplateGroup failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *SystemTemplateGroup) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var SystemTemplateGroupList []*SystemTemplateGroup
	var ResultData Result
	var err error
	cond := orm.NewCondition()

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if this.Account != "" {
		cond = cond.And("account", this.Account)
	}
	if this.Name != "" {
		cond = cond.And("name__contains", this.Name)
	}
	if this.Type != "" {
		cond = cond.And("type", this.Type)
	}
	if this.Version != All {
		cond = cond.And("name__contains", this.Name)
	}
	if this.Status != "" && this.Status != All {
		cond = cond.And("status", this.Status)
	}

	_, err = o.QueryTable(utils.SYSTemplateGroup).SetCond(cond).RelatedSel().Limit(limit, from).All(&SystemTemplateGroupList)
	for _, systemTemplateGroup := range SystemTemplateGroupList {
		o.LoadRelated(systemTemplateGroup, "SystemTemplate")
		o.LoadRelated(systemTemplateGroup, "Job")
	}
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetSYSTemplateGroupErr
		logs.Error("Get SystemTemplateGroup List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.SYSTemplateGroup).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = SystemTemplateGroupList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *SystemTemplateGroup) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	_, err := o.QueryTable(utils.SYSTemplateGroup).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteSYSTemplateGroupErr
		logs.Error("Delete SystemTemplateGroup failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}

func (this *SystemTemplateGroup) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditSYSTemplateGroupErr
		logs.Error("Update SystemTemplateGroup: %s failed, code: %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}
