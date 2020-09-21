package job

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
	jobservice "github.com/xiliangMa/diss-backend/service/job"
	"net/http"
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
	job := new(models.Job)
	job.Id = id

	jobService := new(jobservice.JobService)
	jobService.JobParam = job
	result := jobService.CheckRuningTask()
	if result.Code != http.StatusOK {
		this.Data["json"] = result
		this.ServeJSON(false)
		return
	}
	result = jobService.RemoveAssocTasks()
	result = job.Delete()

	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title DeactiveJob
// @Description Deactive Job
// @Param token header string true "authToken"
// @Param id path string "" true "id"
// @Success 200 {object} models.Result
// @router /:id/deactive [put]
func (this *JobController) DeactiveJob() {
	id := this.GetString(":id")
	job := new(models.Job)
	job.Id = id

	jobService := new(jobservice.JobService)
	jobService.JobParam = job
	// 检查是否有在运行的task，如果有运行的直接返回错误
	result := jobService.CheckRuningTask()
	if result.Code != http.StatusOK {
		this.Data["json"] = result
		this.ServeJSON(false)
		return
	}
	// 设置Job的状态为Deactiving，确定其下的task删除完成后再设置Deactived
	jobService.JobParam.Status = models.Job_Status_Deactiving
	result = jobService.ChangeStatus()
	if result.Code != http.StatusOK {
		this.Data["json"] = result
		this.ServeJSON(false)
		return
	}

	// 下发删除Job下的关联的task
	result = jobService.DeactiveTasks()
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
// @Param account query string "admin" false "租户"
// @Success 200 {object} models.Result
// @router /:id/active [put]
func (this *JobController) ActiveJob() {
	id := this.GetString(":id")
	Job := new(models.Job)
	json.Unmarshal(this.Ctx.Input.RequestBody, &Job)
	Job.Id = id

	jobService := jobservice.JobService{}
	jobService.JobParam = Job

	this.Data["json"] = jobService.ActiveJob()
	this.ServeJSON(false)
}
