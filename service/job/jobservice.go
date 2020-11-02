package job

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/securitycheck"
	taskservice "github.com/xiliangMa/diss-backend/service/task"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type JobService struct {
	JobParam          *models.Job
	securityCheckList []*models.SecurityCheck
	TaskStatusFilter  string
}

func (this *JobService) ActiveJob() models.Result {
	var ResultData models.Result
	job := this.JobParam.Get()

	if job != nil && job.Id != "" {
		// job 处于禁用状态，不执行生成，返回错误信息
		if job.Status == models.Job_Status_Deactiving {
			ResultData.Code = utils.JobIsDisabledErr
			msg := fmt.Sprintf("Id: %s, Job: %s is Deactiving , pls wait it finish or deactive again", job.Id, job.Name)
			logs.Error(msg)
			ResultData.Message = msg
			return ResultData
		}

		// 通过job内容，先将json配置字段清空后 构建checklist
		job.SystemTemplate.CheckControlPlaneJson = ""
		job.SystemTemplate.CheckEtcdJson = ""
		job.SystemTemplate.CheckManagedServicesJson = ""
		job.SystemTemplate.CheckMasterJson = ""
		job.SystemTemplate.CheckNodeJson = ""
		job.SystemTemplate.CheckPoliciesJson = ""
		// 检测主机资源是否已经授权，如果没有授权返回错误
		if len(job.HostConfig) > 0 {
			for _, host := range job.HostConfig {
				if !host.IsLicensed {
					ResultData.Code = utils.LicenseHostErr
					ResultData.Message = "LicenseHostErr"
					return ResultData
				}
			}
		}
		this.securityCheckList = this.BuildCheckList(job)

		if len(this.securityCheckList) > 0 {
			// 通过checklist生成task
			securityCheckService := this.InitsecurityCheckService()
			// 下发task，更新job状态
			ResultData = securityCheckService.DeliverTask()
			if ResultData.Code == utils.LicenseHostErr {
				return ResultData
			}
			job.Status = models.Job_Status_Active
			job.IsUpdateHost = false
			job.Update()
		}
	}

	return ResultData
}

func (this *JobService) InitsecurityCheckService() securitycheck.SecurityCheckService {

	batch := time.Now().Unix()
	SecCheckListModel := models.SecurityCheckList{CheckList: this.securityCheckList}
	securityCheckService := securitycheck.SecurityCheckService{SecurityCheckList: &SecCheckListModel, Batch: batch}

	return securityCheckService
}

func (this *JobService) BuildCheckList(job *models.Job) []*models.SecurityCheck {
	targetType := job.SystemTemplate.Type //目标类型为检查模板的类型
	resType := models.SC_Type_Host        // 资源类型包括 主机 容器 镜像
	var secCheckList []*models.SecurityCheck
	if (targetType == models.TMP_Type_IDS_Docker) || (targetType == models.TMP_Type_DockerVS) {
		resType = models.Sc_Type_Container
	}
	if (resType == models.SC_Type_Host) && (job.HostConfig != nil) {
		for _, checkHost := range job.HostConfig {
			seccheck := this.SetTPLType(targetType, resType)
			seccheck.Host = checkHost
			seccheck.Job = job
			secCheckList = append(secCheckList, seccheck)
		}
	} else if (resType == models.Sc_Type_Container) && (job.ContainerConfig != nil) {
		for _, checkContainer := range job.ContainerConfig {
			seccheck := this.SetTPLType(targetType, resType)
			seccheck.Container = checkContainer
			seccheck.Job = job
			secCheckList = append(secCheckList, seccheck)
		}
	}
	return secCheckList
}

func (this *JobService) SetTPLType(targetType, resType string) (seccheck *models.SecurityCheck) {
	seccheck = new(models.SecurityCheck)
	seccheck.Type = resType
	switch targetType {
	case models.TMP_Type_BM_Docker, models.TMP_Type_BM_K8S:
		seccheck.DockerCIS = true
	case models.TMP_Type_DockerVS, models.TMP_Type_HostVS:
		seccheck.VirusScan = true
	case models.TMP_Type_LS:
		seccheck.LeakScan = true
	}
	return seccheck
}

//判断通过此job生成的task中是否有运行的
func (this *JobService) CheckRuningTask() models.Result {
	var result models.Result

	this.TaskStatusFilter = models.Task_Status_Running
	runingTasks := this.GetTasksOfJobByStatus()

	if len(runingTasks) > 0 {
		result.Code = utils.TaskIsRunningWarn
		msg := fmt.Sprintf("Some Task is Running, code: %d", result.Code)
		logs.Warn(msg)
		result.Message = msg
		result.Data = runingTasks
	}
	result.Code = http.StatusOK
	return result
}

func (this *JobService) GetTasksOfJobByStatus() []*models.Task {
	var tasks []*models.Task
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	job := this.JobParam.Get()
	if job.Id != "" {
		o.LoadRelated(job, "Task")
		for _, task := range job.Task {
			if task.Status == this.TaskStatusFilter {
				tasks = append(tasks, task)
			}
		}
	}
	return tasks
}

func (this *JobService) RemoveAssocTasks() models.Result {
	var result models.Result
	o := orm.NewOrm()
	err := o.Begin()
	//移除相关的task
	job := this.JobParam
	if job.Id != "" {
		o.LoadRelated(job, "Task")
		for _, task := range job.Task {
			taskService := new(taskservice.TaskService)
			taskService.Task = task
			result = taskService.RemoveTask()
			if result.Code != http.StatusOK {
				o.Rollback()
				result.Code = utils.DeleteTaskErr
				msg := fmt.Sprintf("Delete assoc task error, Id: %s , code: %d", task.Id, result.Code)
				logs.Error(msg)
				result.Message = msg
				return result
			}
		}
	}
	err = o.Commit()
	if err != nil {
		result.Code = utils.TaskCommitErr
		result.Message = err.Error()
		logs.Error("Commit Job: %s failed, code: %d, err: %s", job.Name, result.Code, result.Message)
		return result
	}
	result.Code = http.StatusOK
	return result
}

func (this *JobService) DeactiveTasks() models.Result {
	var result models.Result

	// 删除Pending和Unavailable状态task
	taskCount := 0
	job := this.JobParam
	if job.Id != "" {
		jobObj := job.Get()
		for _, task := range jobObj.Task {
			taskService := new(taskservice.TaskService)
			taskService.Task = task

			if task.Status == models.Task_Status_Pending || task.Status == models.Task_Status_Unavailable {
				taskCount++
				task.Action = models.Task_Action_Deactive
				result = taskService.RemoveTask()
			}
		}
	}

	result.Code = http.StatusOK
	result.Data = taskCount
	return result
}

func (this *JobService) ChangeStatus() models.Result {
	result := models.Result{Code: http.StatusOK}

	job := this.JobParam
	jobObj := job.Get()
	jobObj.Status = this.JobParam.Status
	jobObj.IsUpdateHost = false
	result = jobObj.Update()

	return result
}
