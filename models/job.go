package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type Job struct {
	Id                  string               `orm:"pk;" description:"(job id)"`
	Account             string               `orm:"default(admin)" description:"(租户)"`
	Name                string               `orm:"" description:"(名称)"`
	Description         string               `orm:"" description:"(描述)"`
	Spec                string               `orm:"" description:"(定时器)"`
	Type                string               `orm:"" description:"(类型 重复执行 单次执行)"`
	SystemTemplateType  string               `orm:"default(Disable)" description:"(模板类型)"`
	Status              string               `orm:"default(Disable)" description:"(状态: 执行中、启用、禁用)"`
	CreateTime          time.Time            `orm:"auto_now_add;type(datetime)" description:"(创建时间)"`
	UpdateTime          time.Time            `orm:"null;auto_now;type(datetime)" description:"(更新时间)"`
	SystemTemplate      *SystemTemplate      `orm:"rel(fk);null;" description:"(策略)"`
	SystemTemplateGroup *SystemTemplateGroup `orm:"rel(fk);null;" description:"(策略组)"`
	Task                []*Task              `orm:"reverse(many);null" description:"(任务列表)"`
	HostConfig          []*HostConfig        `orm:"reverse(many);null" description:"(主机列表)"`
	ContainerConfig     []*ContainerConfig   `orm:"reverse(many);null" description:"(容器列表)"`
	JobLevel            string               `orm:"default(System)" description:"(任务级别)"`
	IsUpdateHost        bool                 `orm:"-" description:"(是否联动更新相关host信息)"`
}

type JobInterface interface {
	Add() Result
	List(from, limit int) Result
	Delete() Result
	Get(id string) Job
	GetDefaultJob() map[string]*Job
}

func (this *Job) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var JobList []*Job
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
		cond = cond.And("name__icontains", this.Name)
	}
	if this.JobLevel == "" {
		cond = cond.And("job_level", Job_Level_User)
	}
	if this.HostConfig != nil && len(this.HostConfig) > 0 {
		hostid := this.HostConfig[0].Id
		cond = cond.And("HostConfig__HostConfig__Id", hostid)
	}
	_, err = o.QueryTable(utils.Job).SetCond(cond).RelatedSel().Limit(limit, from).OrderBy("-create_time").All(&JobList)

	for _, job := range JobList {
		o.LoadRelated(job, "HostConfig")
		o.LoadRelated(job, "ContainerConfig")
		o.LoadRelated(job, "Task")

		if job.Task != nil {
			for _, task := range job.Task {
				if task.Host != nil {
					task.Host = task.Host.Get()
				}
			}
		}
	}

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetJobErr
		logs.Error("Get Job List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.Job).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = JobList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *Job) Get() *Job {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	job := new(Job)
	var ResultData Result
	var err error
	cond := orm.NewCondition()

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	err = o.QueryTable(utils.Job).SetCond(cond).RelatedSel().One(job)

	if job.Id != "" {
		o.LoadRelated(job, "HostConfig")
		o.LoadRelated(job, "ContainerConfig")
		o.LoadRelated(job, "Task")
	}
	if err != nil && utils.IgnoreQuerySeterNoRowFoundErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetJobErr
		logs.Error("Get Job Item failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return job
	}

	return job
}

func (this *Job) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	err := o.Begin()

	uid, _ := uuid.NewV4()
	this.Id = uid.String()

	_, err = o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		o.Rollback()
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddJobErr
		logs.Error("Add Job failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	result := this.fillRelationData(o)
	if result.Code != 0 && result.Code != utils.ContainerExistInJobErr && result.Code != utils.HostExistInJobErr {
		return result
	}

	err = o.Commit()
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Code = utils.JobCommitErr
		ResultData.Message = err.Error()
		logs.Error("Commit Job: %s failed, code: %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Job) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	err := o.Begin()
	// 移除job-host  , job-container关联
	m2mhost := o.QueryM2M(this, "HostConfig")
	_, err = m2mhost.Clear()
	m2mcont := o.QueryM2M(this, "ContainerConfig")
	_, err = m2mcont.Clear()
	if err != nil {
		o.Rollback()
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteJobErr
		logs.Error("Delete Job: %s failed, code: %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
		return ResultData
	}
	// 移除job本身
	cond := orm.NewCondition()

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}

	_, err = o.QueryTable(utils.Job).SetCond(cond).Delete()
	if err != nil {
		o.Rollback()
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteJobErr
		logs.Error("Delete Job failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	err = o.Commit()
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Code = utils.JobCommitErr
		ResultData.Message = err.Error()
		logs.Error("Commit Job: %s failed, code: %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	return ResultData
}

func (this *Job) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	err := o.Begin()
	_, err = o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditJobErr
		o.Rollback()
		logs.Error("Update Job: %s failed, code: %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
		return ResultData
	}

	if this.IsUpdateHost {
		result := this.fillRelationData(o)
		if result.Code != 0 && result.Code != utils.ContainerExistInJobErr && result.Code != utils.HostExistInJobErr {
			return result
		}
	}

	err = o.Commit()
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Code = utils.JobCommitErr
		ResultData.Message = err.Error()
		logs.Error("Commit Job: %s failed, code: %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Job) fillRelationData(o orm.Ormer) Result {
	var ResultData Result
	m2m := o.QueryM2M(this, "HostConfig")
	m2m.Clear()
	for _, hostconfig := range this.HostConfig {
		if m2m.Exist(hostconfig) != true {
			_, err := m2m.Add(hostconfig)
			if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
				ResultData.Message = err.Error()
				ResultData.Code = utils.RelationJobHostErr
				o.Rollback()
				logs.Error("Relation Job %s to HostConifg error, code %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
				return ResultData
			}
		} else {
			ResultData.Message = "Host Resource Exist"
			ResultData.Code = utils.HostExistInJobErr
			logs.Warn("Relation Job %s to HostConifg error, code %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
		}
	}
	m2m = o.QueryM2M(this, "ContainerConfig")
	m2m.Clear()
	for _, containerconfig := range this.ContainerConfig {
		if m2m.Exist(containerconfig) != true {
			_, err := m2m.Add(containerconfig)
			if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
				ResultData.Message = err.Error()
				ResultData.Code = utils.RelationJobContainerErr
				o.Rollback()
				logs.Error("Relation Job %s to ContainerConfig error, code %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
				return ResultData
			}
		} else {
			ResultData.Message = "Container Resource Exist"
			ResultData.Code = utils.ContainerExistInJobErr
			logs.Warn("Relation Job %s to HostConifg error, code %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
		}
	}
	return ResultData
}

func (this *Job) GetDefaultJob() map[string]*Job {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var systemJobList []*Job
	defaultJobList := make(map[string]*Job)
	var err error
	cond := orm.NewCondition()
	cond = cond.And("job_level", Job_Level_System)
	_, err = o.QueryTable(utils.Job).SetCond(cond).RelatedSel().All(&systemJobList)
	if err != nil {
		logs.Error("Get SystemJob List failed, code: %d, err: %s", utils.GetSYSJobErr, err.Error())
		return nil
	}

	for _, systemJob := range systemJobList {
		defaultJobList[systemJob.SystemTemplateType] = systemJob
	}
	return defaultJobList
}
