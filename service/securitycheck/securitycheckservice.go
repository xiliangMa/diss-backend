package securitycheck

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	uuid "github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/ws"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type SecurityCheckService struct {
	*models.SecurityCheckList
	DefaultTMP           map[string]*models.SystemTemplate
	DefaultJob           map[string]*models.Job
	Batch                int64
	CurrentBatchTaskList []*models.Task
}

func (this *SecurityCheckService) PrePare() {
	logs.Info("################ PrePare work <<<start>>> ################")
	this.PrePareDefaultJob()
	this.PrePareDefaultTMP()
	logs.Info("PrePare task start ......")
	for _, securityCheck := range this.SecurityCheckList.CheckList {
		switch securityCheck.Type {
		case models.SC_Type_Host:
			this.PrePareTask(securityCheck)
		case models.Sc_Type_Container:
			this.PrePareTask(securityCheck)
		}
	}
	logs.Info("PrePare task end ......")

	this.GetCurrentBatchTask()

	logs.Info("################ PrePare work <<<end>>> ################")
}

func (this *SecurityCheckService) PrePareTask(securityCheck *models.SecurityCheck) {

	//默认Job
	Job_Type_BM_Docker := this.DefaultJob[models.TMP_Type_BM_Docker]
	Job_Type_BM_K8S := this.DefaultJob[models.TMP_Type_BM_K8S]
	Job_Type_DockerVS := this.DefaultJob[models.TMP_Type_DockerVS]
	Job_Type_HostVS := this.DefaultJob[models.TMP_Type_HostVS]

	// 默认模板
	TMP_Type_BM_Docker_DT := this.DefaultTMP[models.TMP_Type_BM_Docker]
	TMP_Type_BM_K8S_DT := this.DefaultTMP[models.TMP_Type_BM_K8S]
	TMP_Type_DockerVS_DT := this.DefaultTMP[models.TMP_Type_DockerVS]
	TMP_Type_HostVS_DT := this.DefaultTMP[models.TMP_Type_HostVS]
	//TMP_Type_LS_DT := defaultTMP[models.TMP_Type_LS]

	// 检测是否是系统Job
	if securityCheck.Job == nil {
		if securityCheck.DockerBenchMarkCheck {
			securityCheck.Job = Job_Type_BM_Docker
		}
		if securityCheck.KubeBenchMarkCheck {
			securityCheck.Job = Job_Type_BM_K8S
		}

		//todo
		if securityCheck.VirusScan && securityCheck.Type == models.SC_Type_Host {
			securityCheck.Job = Job_Type_HostVS
		}

		if securityCheck.VirusScan && securityCheck.Type == models.Sc_Type_Container {
			securityCheck.Job = Job_Type_DockerVS
		}
		//todo
		//if securityCheck.LeakScan {
		//}
	}

	taskpre := "系统任务-"
	if securityCheck.Job.JobLevel == models.Job_Level_User {
		taskpre = "用户任务-"
	}

	if securityCheck.DockerBenchMarkCheck {

		if securityCheck.Job.JobLevel == models.Job_Level_System || securityCheck.Job.SystemTemplate.Type == models.TMP_Type_BM_Docker {
			this.genBenchmarkTask(securityCheck, TMP_Type_BM_Docker_DT, taskpre)
		}

		if securityCheck.Job.JobLevel == models.Job_Level_System || securityCheck.Job.SystemTemplate.Type == models.TMP_Type_BM_K8S {
			this.genBenchmarkTask(securityCheck, TMP_Type_BM_K8S_DT, taskpre)
		}

	}
	if securityCheck.VirusScan {
		//病毒
		task := new(models.Task)
		uid, _ := uuid.NewV4()
		task.Id = uid.String()
		if securityCheck.Type != models.SC_Type_Host {
			//容器病毒
			logs.Info("PrePare task, Type:  %s , Task Id: %s ......", models.TMP_Type_DockerVS, uid)
			task.SystemTemplate = TMP_Type_DockerVS_DT
			task.Description = taskpre + models.TMP_Type_DockerVS
		} else {
			//主机病毒
			logs.Info("PrePare task, Type:  %s , Task Id: %s ......", models.TMP_Type_HostVS, uid)
			task.SystemTemplate = TMP_Type_HostVS_DT
			task.Description = taskpre + models.TMP_Type_HostVS
		}

		if (*securityCheck).Job.Type == "" {
			task.Type = models.Job_Type_Once
		} else {
			task.Type = securityCheck.Job.Type
		}
		task.Name = taskpre + task.Id
		task.Host = securityCheck.Host
		task.Container = securityCheck.Container
		task.Batch = this.Batch
		task.Status = models.Task_Status_Pending
		task.Account = securityCheck.Job.Account
		task.Job = securityCheck.Job
		task.Spec = securityCheck.Job.Spec
		//添加task记录
		task.Add()

		//添加任务日志
		taskLog := models.TaskLog{}
		taskRawInfo, _ := json.Marshal(task)
		taskLog.Task = string(taskRawInfo)
		taskLog.Account = securityCheck.Job.Account
		taskLog.Level = models.Log_level_Info
		taskLog.RawLog = fmt.Sprintf("Add security check task, Id: %s, Type: %s, Batch: %v, Status: %s",
			task.Id, task.Type, task.Batch, task.Status)
		taskLog.Add()
	}
	if securityCheck.LeakScan {
		//漏洞
		task := new(models.Task)
		uid, _ := uuid.NewV4()
		task.Id = uid.String()
		logs.Info("PrePare task, Type:  %s , Task Id: %s ......", models.TMP_Type_LS, uid)
		//task.SystemTemplate = TMP_Type_LS_DT
		task.Name = taskpre + task.Id
		task.Description = taskpre + models.TMP_Type_LS
		task.Host = securityCheck.Host
		task.Batch = this.Batch
		task.Account = securityCheck.Job.Account
		task.Job = securityCheck.Job
		task.Spec = securityCheck.Job.Spec
		//添加task记录
		task.Add()
	}
}

func (this *SecurityCheckService) genBenchmarkTask(securityCheck *models.SecurityCheck, Default_Template *models.SystemTemplate, taskpre string) *models.Task {
	task := new(models.Task)
	uid, _ := uuid.NewV4()
	task.Id = uid.String()
	if securityCheck.Job.SystemTemplate.Type == models.TMP_Type_BM_Docker {
		//基线-Docker
		logs.Info("PrePare task, Type:  %s , Task Id: %s ......", models.TMP_Type_BM_Docker, uid)
		task.Description = taskpre + models.TMP_Type_BM_Docker
	} else {
		//基线-K8S
		logs.Info("PrePare task, Type:  %s , Task Id: %s ......", models.TMP_Type_BM_K8S, uid)
		task.Description = taskpre + models.TMP_Type_BM_K8S
	}
	task.Type = securityCheck.Job.Type
	task.SystemTemplate = Default_Template
	task.Name = taskpre + task.Id

	task.Host = securityCheck.Host
	task.Batch = this.Batch
	task.Status = models.Task_Status_Pending
	task.Account = securityCheck.Job.Account
	task.Job = securityCheck.Job
	task.Spec = securityCheck.Job.Spec
	task.Add()
	taskLog := models.TaskLog{}
	taskRawInfo, _ := json.Marshal(task)
	taskLog.Task = string(taskRawInfo)
	taskLog.Account = securityCheck.Job.Account
	taskLog.Level = models.Log_level_Info
	taskLog.RawLog = fmt.Sprintf("Add security check task, Id: %s, Type: %s, Batch: %v, Status: %s",
		task.Id, task.Type, task.Batch, task.Status)
	return task
}

func (this *SecurityCheckService) PrePareDefaultTMP() map[string]*models.SystemTemplate {
	if this.DefaultTMP == nil {
		st := new(models.SystemTemplate)
		this.DefaultTMP = st.GetDefaultTemplate()
	}
	logs.Info("PrePare Default Template: %v", this.DefaultTMP)
	return this.DefaultTMP
}

func (this *SecurityCheckService) PrePareDefaultJob() map[string]*models.Job {
	if this.DefaultTMP == nil {
		defaultJob := new(models.Job)
		this.DefaultJob = defaultJob.GetDefaultJob()
	}
	logs.Info("PrePare Default Job: %v", this.DefaultJob)
	return this.DefaultJob
}

func (this *SecurityCheckService) GetCurrentBatchTask() []*models.Task {
	if this.CurrentBatchTaskList == nil {
		task := new(models.Task)
		task.Batch = this.Batch
		if err, taskList := task.GetCurrentBatchTaskList(); err == nil {
			this.CurrentBatchTaskList = taskList
		}
	}
	logs.Info("Get Current Batch Task List: %v", this.CurrentBatchTaskList)
	return this.CurrentBatchTaskList
}

func (this *SecurityCheckService) DeliverTask(isPrepare bool) models.Result {
	var ResultData models.Result

	if isPrepare {
		this.PrePare()
	}

	wsDelive := ws.WSDeliverService{
		Hub:                  models.WSHub,
		Bath:                 this.Batch,
		CurrentBatchTaskList: this.CurrentBatchTaskList,
	}
	if utils.IsEnableNats() {
		go wsDelive.DeliverTaskToNats()
	} else {
		go wsDelive.DeliverTask()

	}

	ResultData.Code = http.StatusOK
	data := make(map[string]interface{})
	data["items"] = this.CurrentBatchTaskList
	ResultData.Data = data

	return ResultData
}
