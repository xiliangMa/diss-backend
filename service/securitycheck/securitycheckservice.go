package securitycheck

import (
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
	Batch                int64
	CurrentBatchTaskList []*models.Task
	Account              string
}

func (this *SecurityCheckService) PrePare() {
	logs.Info("################ PrePare work <<<start>>> ################")
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
	TMP_Type_BM_Docker_DT := this.DefaultTMP[models.TMP_Type_BM_Docker]
	TMP_Type_BM_K8S_DT := this.DefaultTMP[models.TMP_Type_BM_K8S]
	TMP_Type_DockerVS_DT := this.DefaultTMP[models.TMP_Type_DockerVS]
	TMP_Type_HostVS_DT := this.DefaultTMP[models.TMP_Type_HostVS]
	//TMP_Type_LS_DT := defaultTMP[models.TMP_Type_LS]
	if securityCheck.BenchMarkCheck {
		dockerTask := new(models.Task)
		uid1, _ := uuid.NewV4()
		dockerTask.Id = uid1.String()
		logs.Info("PrePare task, Type:  %s , Task Id: %s ......", models.TMP_Type_BM_Docker, uid1)

		//基线-Docker
		dockerTask.Type = models.Job_Type_Once
		dockerTask.SystemTemplate = TMP_Type_BM_Docker_DT
		dockerTask.Name = "System-Task-" + dockerTask.Id
		dockerTask.Description = "System-Task-" + models.TMP_Type_BM_Docker
		dockerTask.Host = securityCheck.Host
		dockerTask.Batch = this.Batch
		dockerTask.Status = models.Task_Status_Pending
		dockerTask.Account = this.Account

		k8sTask := new(models.Task)
		uid2, _ := uuid.NewV4()
		k8sTask.Id = uid2.String()
		logs.Info("PrePare task, Type:  %s , Task Id: %s ......", models.TMP_Type_BM_K8S, uid2)
		//基线-K8S
		k8sTask.Type = models.Job_Type_Once
		k8sTask.SystemTemplate = TMP_Type_BM_K8S_DT
		k8sTask.Name = "System-Task-" + k8sTask.Id
		k8sTask.Description = "System-Task-" + models.TMP_Type_BM_K8S
		k8sTask.Host = securityCheck.Host
		k8sTask.Batch = this.Batch
		k8sTask.Status = models.Task_Status_Pending
		k8sTask.Account = this.Account

		//添加task记录
		dockerTask.Add()
		k8sTask.Add()

		//添加任务日志
		dockerTaskLog := models.TaskLog{}
		dockerTaskLog.Task = dockerTask
		dockerTaskLog.Account = this.Account
		dockerTaskLog.Level = models.Log_level_Info
		dockerTaskLog.RawLog = fmt.Sprintf("Add security check task, Id: %s, Type: %s, Batch: %v, Status: %s",
			dockerTask.Id, dockerTask.Type, dockerTask.Batch, dockerTask.Status)
		k8sTaskLog := models.TaskLog{}
		k8sTaskLog.Task = dockerTask
		k8sTaskLog.Account = this.Account
		k8sTaskLog.Level = models.Log_level_Info
		k8sTaskLog.RawLog = fmt.Sprintf("Add security check task, Id: %s, Type: %s, Batch: %v, Status: %s",
			k8sTask.Id, k8sTask.Type, k8sTask.Batch, k8sTask.Status)
		dockerTaskLog.Add()
		k8sTaskLog.Add()
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
			task.Description = "System-Task-" + models.TMP_Type_DockerVS
		} else {
			//主机病毒
			logs.Info("PrePare task, Type:  %s , Task Id: %s ......", models.TMP_Type_HostVS, uid)
			task.SystemTemplate = TMP_Type_HostVS_DT
			task.Description = "System-Task-" + models.TMP_Type_HostVS
		}
		task.Type = models.Job_Type_Once
		task.Name = "System-Task-" + task.Id
		task.Host = securityCheck.Host
		task.Container = securityCheck.Container
		task.Batch = this.Batch
		task.Status = models.Task_Status_Pending
		task.Account = this.Account

		//添加task记录
		task.Add()

		//添加任务日志
		taskLog := models.TaskLog{}
		taskLog.Task = task
		taskLog.Account = this.Account
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
		task.Name = "System-Task-" + task.Id
		task.Description = "System-Task-" + models.TMP_Type_LS
		task.Host = securityCheck.Host
		task.Batch = this.Batch
		//添加task记录
		task.Add()
	}
}

func (this *SecurityCheckService) PrePareDefaultTMP() map[string]*models.SystemTemplate {
	if this.DefaultTMP == nil {
		st := new(models.SystemTemplate)
		this.DefaultTMP = st.GetDefaultTemplate()
	}
	logs.Info("PrePare Default Template: %v", this.DefaultTMP)
	return this.DefaultTMP
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

func (this *SecurityCheckService) DeliverTask() models.Result {
	var ResultData models.Result
	this.PrePare()
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
