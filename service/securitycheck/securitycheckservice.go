package securitycheck

import (
	"github.com/astaxie/beego/logs"
	uuid "github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/models/bean"
	"github.com/xiliangMa/diss-backend/models/global"
	"github.com/xiliangMa/diss-backend/models/job"
	msecuritypolicy "github.com/xiliangMa/diss-backend/models/securitypolicy"
	"github.com/xiliangMa/diss-backend/service/ws"
	"net/http"
)

type SecurityCheckService struct {
	*bean.SecurityCheckList
	DefaultTMP           map[string]*msecuritypolicy.SystemTemplate
	Bath                 int64
	CurrentBatchTaskList []*job.Task
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

func (this *SecurityCheckService) PrePareTask(securityCheck *bean.SecurityCheck) {
	TMP_Type_BM_Docker_DT := this.DefaultTMP[models.TMP_Type_BM_Docker]
	TMP_Type_BM_K8S_DT := this.DefaultTMP[models.TMP_Type_BM_K8S]
	//TMP_Type_VS_DT := defaultTMP[models.TMP_Type_VS]
	//TMP_Type_LS_DT := defaultTMP[models.TMP_Type_LS]
	if securityCheck.BenchMarkCheck {
		dockerTask := new(job.Task)
		uid1, _ := uuid.NewV4()
		dockerTask.Id = uid1.String()
		logs.Info("PrePare task, Type:  %s , Task Id: %s ......", models.TMP_Type_BM_Docker, uid1)

		//基线-Docker
		dockerTask.Type = models.Job_Type_Once
		dockerTask.SystemTemplate = TMP_Type_BM_Docker_DT
		dockerTask.Name = "System-Task-" + dockerTask.Id
		dockerTask.Description = "System-Task-" + models.TMP_Type_BM_Docker
		dockerTask.Host = securityCheck.Host
		dockerTask.Batch = this.Bath
		dockerTask.Status = models.Task_Status_Pending

		k8sTask := new(job.Task)
		uid2, _ := uuid.NewV4()
		k8sTask.Id = uid2.String()
		logs.Info("PrePare task, Type:  %s , Task Id: %s ......", models.TMP_Type_BM_K8S, uid2)
		//基线-K8S
		k8sTask.Type = models.Job_Type_Once
		k8sTask.SystemTemplate = TMP_Type_BM_K8S_DT
		k8sTask.Name = "System-Task-" + k8sTask.Id
		k8sTask.Description = "System-Task-" + models.TMP_Type_BM_K8S
		k8sTask.Host = securityCheck.Host
		k8sTask.Batch = this.Bath
		k8sTask.Status = models.Task_Status_Pending

		//添加task记录
		dockerTask.Add()
		k8sTask.Add()
	}
	if securityCheck.VirusScan {
		//病毒
		task := new(job.Task)
		uid, _ := uuid.NewV4()
		task.Id = uid.String()
		logs.Info("PrePare task, Type:  %s , Task Id: %s ......", models.TMP_Type_VS, uid)
		//task.SystemTemplate = TMP_Type_VS_DT
		task.Name = "System-Task-" + task.Id
		task.Description = "System-Task-" + models.TMP_Type_VS
		task.Host = securityCheck.Host
		task.Batch = this.Bath
		//添加task记录
		task.Add()
	}
	if securityCheck.LeakScan {
		//漏洞
		task := new(job.Task)
		uid, _ := uuid.NewV4()
		task.Id = uid.String()
		logs.Info("PrePare task, Type:  %s , Task Id: %s ......", models.TMP_Type_LS, uid)
		//task.SystemTemplate = TMP_Type_LS_DT
		task.Name = "System-Task-" + task.Id
		task.Description = "System-Task-" + models.TMP_Type_LS
		task.Host = securityCheck.Host
		task.Batch = this.Bath
		//添加task记录
		task.Add()
	}
}

func (this *SecurityCheckService) PrePareDefaultTMP() map[string]*msecuritypolicy.SystemTemplate {
	if this.DefaultTMP == nil {
		st := new(msecuritypolicy.SystemTemplate)
		this.DefaultTMP = st.GetDefaultTemplate()
	}
	logs.Info("PrePare Default Template: %v", this.DefaultTMP)
	return this.DefaultTMP
}

func (this *SecurityCheckService) GetCurrentBatchTask() []*job.Task {
	if this.CurrentBatchTaskList == nil {
		task := new(job.Task)
		task.Batch = this.Bath
		if err, taskList := task.GetCurrentBatchTaskList(); err == nil {
			this.CurrentBatchTaskList = taskList
		}
	}
	logs.Info("Get Current Batch Task List: %v", this.CurrentBatchTaskList)
	return this.CurrentBatchTaskList
}

func (this *SecurityCheckService) DeliverTask(isNats bool) models.Result {
	var ResultData models.Result
	this.PrePare()
	wsDelive := ws.WSDeliverService{
		Hub:                  global.WSHub,
		Bath:                 this.Bath,
		CurrentBatchTaskList: this.CurrentBatchTaskList,
	}
	if isNats {
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
