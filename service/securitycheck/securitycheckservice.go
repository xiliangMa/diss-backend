package securitycheck

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/ws"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type SecurityCheckService struct {
	*models.SecurityCheckList
	ClusterCheckObject   *models.ClusterCheck
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
			// 初始化 Host
			GetCheckHost(securityCheck)
			job := securityCheck.Job
			if securityCheck.KubenetesCIS {
				securityCheck := models.SecurityCheck{
					KubenetesCIS: securityCheck.KubenetesCIS,
					Host:         securityCheck.Host,
					Type:         models.SC_Type_Host,
					Job:          job,
				}
				this.PrePareTask(&securityCheck)
			}
			if securityCheck.DockerCIS {
				securityCheck := models.SecurityCheck{
					DockerCIS: securityCheck.DockerCIS,
					Host:      securityCheck.Host,
					Type:      models.SC_Type_Host,
					Job:       job,
				}
				this.PrePareTask(&securityCheck)
			}
			if securityCheck.VirusScan {
				securityCheck := models.SecurityCheck{
					VirusScan: securityCheck.VirusScan,
					Host:      securityCheck.Host,
					Type:      models.SC_Type_Host,
					Job:       job,
				}
				this.PrePareTask(&securityCheck)
			}
			if securityCheck.LeakScan {
				securityCheck := models.SecurityCheck{
					LeakScan: securityCheck.LeakScan,
					Host:     securityCheck.Host,
					Type:     models.SC_Type_Host,
					Job:      job,
				}
				this.PrePareTask(&securityCheck)
			}

		case models.Sc_Type_Container:
			GetCheckContainer(securityCheck)
			job := securityCheck.Job
			if securityCheck.KubenetesCIS {
				securityCheck := models.SecurityCheck{
					KubenetesCIS: securityCheck.KubenetesCIS,
					Container:    securityCheck.Container,
					Type:         models.SC_Type_Host,
					Job:          job,
				}
				this.PrePareTask(&securityCheck)
			}
			if securityCheck.DockerCIS {
				securityCheck := models.SecurityCheck{
					DockerCIS: securityCheck.DockerCIS,
					Container: securityCheck.Container,
					Type:      models.SC_Type_Host,
					Job:       job,
				}
				this.PrePareTask(&securityCheck)
			}
			if securityCheck.VirusScan {
				securityCheck := models.SecurityCheck{
					VirusScan: securityCheck.VirusScan,
					Container: securityCheck.Container,
					Type:      models.SC_Type_Host,
					Job:       job,
				}
				this.PrePareTask(&securityCheck)
			}
			if securityCheck.LeakScan {
				securityCheck := models.SecurityCheck{
					LeakScan:  securityCheck.LeakScan,
					Container: securityCheck.Container,
					Type:      models.SC_Type_Host,
					Job:       job,
				}
				this.PrePareTask(&securityCheck)
			}
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
		if securityCheck.DockerCIS {
			securityCheck.Job = Job_Type_BM_Docker
		}
		if securityCheck.KubenetesCIS {
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

	if securityCheck.Job == nil {
		return
	}
	taskpre := "系统任务-"
	if securityCheck.Job.JobLevel == models.Job_Level_User {
		taskpre = "用户任务-"
	}

	if securityCheck.Job.SystemTemplate.Type == models.TMP_Type_BM_Docker {
		this.genBenchmarkTask(securityCheck, TMP_Type_BM_Docker_DT, taskpre)
	}

	if securityCheck.Job.SystemTemplate.Type == models.TMP_Type_BM_K8S {
		this.genBenchmarkTask(securityCheck, TMP_Type_BM_K8S_DT, taskpre)
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

		// 集群ID
		if this.ClusterCheckObject != nil {
			task.ClusterId = this.ClusterCheckObject.Id
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
		this.saveTaskLog(task, securityCheck)

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
		if this.ClusterCheckObject != nil {
			task.ClusterId = this.ClusterCheckObject.Id
		}
		task.Add()

		//添加task记录
		this.saveTaskLog(task, securityCheck)
	}
}

func GetCheckHost(securityCheck *models.SecurityCheck) (*models.SecurityCheck, error) {
	hostId := securityCheck.Host.Id
	if hostId == "" {
		return nil, errors.Errorf("HostId is not null")
	}

	// 获取主机
	host := models.HostConfig{Id: hostId}
	checkHost := host.Get()
	if checkHost == nil {
		return nil, errors.Errorf("Host not found")
	}
	securityCheck.Host = checkHost
	return securityCheck, nil
}

func GetCheckContainer(securityCheck *models.SecurityCheck) (*models.SecurityCheck, error) {
	containerId := securityCheck.Container.Id
	if containerId == "" {
		return nil, errors.Errorf("ContainerId is not null")
	}

	// 获取主机
	container := models.ContainerConfig{Id: containerId}
	checkContainer := container.Get()
	if checkContainer == nil {
		return nil, errors.Errorf("Container not found")
	}
	securityCheck.Container = checkContainer
	return securityCheck, nil
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
	if this.ClusterCheckObject != nil {
		task.ClusterId = this.ClusterCheckObject.Id
	}

	task.Add()
	this.saveTaskLog(task, securityCheck)

	return task
}

func (this *SecurityCheckService) saveTaskLog(task *models.Task, securityCheck *models.SecurityCheck) {
	taskLog := models.TaskLog{}
	taskRawInfo, _ := json.Marshal(task)
	taskLog.Task = string(taskRawInfo)
	taskLog.Account = securityCheck.Job.Account
	taskLog.Level = models.Log_level_Info
	taskLog.RawLog = fmt.Sprintf("Add security check task, Id: %s, Type: %s, Batch: %v, KStatus: %s",
		task.Id, task.Type, task.Batch, task.Status)
	taskLog.Add()
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

func (this *SecurityCheckService) ClusterCheck() models.Result {
	result := models.Result{Code: http.StatusOK}
	this.PrePareDefaultJob()
	this.PrePareDefaultTMP()

	// 获取集群内主机
	host := models.HostConfig{ClusterId: this.ClusterCheckObject.Id}
	hostResult := host.List(0, 0)
	if hostResult.Code != http.StatusOK {
		result.Code = utils.GetHostListForClusterErr
		return result
	}
	hostList := hostResult.Data.(map[string]interface{})[models.Result_Items].([]*models.HostConfig)

	// 生成task
	for _, host := range hostList {
		if this.ClusterCheckObject.KubenetesCIS {
			securityCheck := models.SecurityCheck{
				KubenetesCIS: this.ClusterCheckObject.KubenetesCIS,
				Host:         host,
				Type:         models.SC_Type_Host,
			}
			this.PrePareTask(&securityCheck)
		}
		if this.ClusterCheckObject.DockerCIS {
			securityCheck := models.SecurityCheck{
				DockerCIS: this.ClusterCheckObject.DockerCIS,
				Host:      host,
				Type:      models.SC_Type_Host,
			}
			this.PrePareTask(&securityCheck)
		}
		if this.ClusterCheckObject.VirusScan {
			securityCheck := models.SecurityCheck{
				VirusScan: this.ClusterCheckObject.VirusScan,
				Host:      host,
				Type:      models.SC_Type_Host,
			}
			this.PrePareTask(&securityCheck)
		}
		if this.ClusterCheckObject.LeakScan {
			securityCheck := models.SecurityCheck{
				LeakScan: this.ClusterCheckObject.LeakScan,
				Host:     host,
				Type:     models.SC_Type_Host,
			}
			this.PrePareTask(&securityCheck)
		}
	}

	// 下发
	wsDelive := ws.WSDeliverService{
		Hub:                  models.WSHub,
		Bath:                 this.Batch,
		CurrentBatchTaskList: this.CurrentBatchTaskList,
	}
	go wsDelive.DeliverTaskToNats()

	data := make(map[string]interface{})
	data[models.Result_Items] = this.ClusterCheckObject
	data[models.Result_Total] = 0

	result.Data = data
	// 返回批次
	return result

}
