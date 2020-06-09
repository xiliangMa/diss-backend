package job

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/job"
	"github.com/xiliangMa/diss-backend/service/securitycheck"
	"time"
)

// Job 接口
type JobController struct {
	beego.Controller
}

// @Title GetJobList
// @Description Get Job List
// @Param token header string true "authToken"
// @Param body body models.Job false "Job"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *JobController) GetJobList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	Job := new(models.Job)
	json.Unmarshal(this.Ctx.Input.RequestBody, &Job)
	this.Data["json"] = Job.List(from, limit)
	this.ServeJSON(false)

}

// @Title DeleteJob
// @Description Delete Job
// @Param token header string true "authToken"
// @Param id path string "" true "id"
// @Success 200 {object} models.Result
// @router /:id [delete]
func (this *JobController) DeleteJob() {
	id := this.GetString(":id")
	Job := new(models.Job)
	Job.Id = id
	result := Job.Delete()
	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title AddJob
// @Description Add Job
// @Param token header string true "authToken"
// @Param body body models.Job false "Job"
// @Success 200 {object} models.Result
// @router /add [post]
func (this *JobController) AddJob() {
	Job := new(models.Job)
	json.Unmarshal(this.Ctx.Input.RequestBody, &Job)
	result := Job.Add()
	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title UpdateJob
// @Description Update Job
// @Param token header string true "authToken"
// @Param id path string "" true "id"
// @Param body body models.Job false "Job"
// @Success 200 {object} models.Result
// @router /:id [put]
func (this *JobController) UpdateJob() {
	id := this.GetString(":id")
	Job := new(models.Job)
	json.Unmarshal(this.Ctx.Input.RequestBody, &Job)
	Job.Id = id
	result := Job.Update()
	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title ActiveJob
// @Description Gen Tasks From Job , and Delivery Tasks
// @Param token header string true "authToken"
// @Param id path string "" true "id"
// @Success 200 {object} models.Result
// @router /active/:id [put]
func (this *JobController) ActiveJob() {
	id := this.GetString(":id")
	Job := new(models.Job)
	json.Unmarshal(this.Ctx.Input.RequestBody, &Job)
	Job.Id = id

	jobservice := job.JobService{}
	jobservice.JobParm = Job

	account := this.GetString("account")
	if account == "" {
		account = models.Account_Admin
	}
	batch := time.Now().Unix()
	actChecklist := jobservice.GetCheckList()
	secCheckList := models.SecurityCheckList{CheckList: actChecklist}
	securityCheckService := securitycheck.SecurityCheckService{SecurityCheckList: &secCheckList, Batch: batch, Account: account}

	this.Data["json"] = securityCheckService.DeliverTask()
	this.ServeJSON(false)
}
