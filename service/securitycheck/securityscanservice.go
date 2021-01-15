package securitycheck

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	uuid "github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/base"
	"github.com/xiliangMa/diss-backend/service/ws"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type SecurityScanService struct {
	*models.SecurityCheckParams
	ClusterCheckObject   *models.ClusterCheck
	DefaultTMP           map[string]*models.SystemTemplate
	DefaultJob           map[string]*models.Job
	HostList             []*models.HostConfig
	ContainerList        []*models.ContainerConfig
	ImageList            []*models.ImageConfig
	Job                  *models.Job
	Batch                int64
	CurrentBatchTaskList []*models.Task
	IsSystem             bool
}

func (this *SecurityScanService) PrePare() {
	logs.Info("################ PrePare work <<<start>>> ################")
	// 获取系统Job、模版信息
	this.PrePareDefaultJob()
	this.PrePareDefaultTMP()

	// 构建检查对象、并生成task
	logs.Info("PrePare task start ......")
	switch this.SecurityCheckParams.Type {
	case models.SC_Type_Host:
		// 清空不需要的json配置字段
		if this.Job != nil {
			this.Job.SystemTemplate.CheckControlPlaneJson = ""
			this.Job.SystemTemplate.CheckEtcdJson = ""
			this.Job.SystemTemplate.CheckManagedServicesJson = ""
			this.Job.SystemTemplate.CheckMasterJson = ""
			this.Job.SystemTemplate.CheckNodeJson = ""
			this.Job.SystemTemplate.CheckPoliciesJson = ""
		}
		for _, host := range this.HostList {
			if this.SecurityCheckParams.KubenetesCIS {
				securityCheck := models.SecurityCheck{
					KubenetesCIS: this.SecurityCheckParams.KubenetesCIS,
					Host:         host,
					Type:         models.SC_Type_Host,
					Job:          this.Job,
				}
				this.PrePareTask(&securityCheck)
			}
			if this.SecurityCheckParams.DockerCIS {
				securityCheck := models.SecurityCheck{
					DockerCIS: this.SecurityCheckParams.DockerCIS,
					Host:      host,
					Type:      models.SC_Type_Host,
					Job:       this.Job,
				}
				this.PrePareTask(&securityCheck)
			}
			if this.SecurityCheckParams.VirusScan {
				securityCheck := models.SecurityCheck{
					VirusScan: this.SecurityCheckParams.VirusScan,
					Host:      host,
					Type:      models.SC_Type_Host,
					Job:       this.Job,
				}
				this.PrePareTask(&securityCheck)
			}
			if this.SecurityCheckParams.LeakScan {
				securityCheck := models.SecurityCheck{
					LeakScan: this.SecurityCheckParams.LeakScan,
					Host:     host,
					Type:     models.SC_Type_Host,
					Job:      this.Job,
				}
				this.PrePareTask(&securityCheck)
			}
		}
	case models.Sc_Type_Container:
		for _, object := range this.ContainerList {
			if this.SecurityCheckParams.VirusScan {
				securityCheck := models.SecurityCheck{
					VirusScan: this.SecurityCheckParams.VirusScan,
					Container: object,
					Type:      models.Sc_Type_Container,
					Job:       this.Job,
				}
				this.PrePareTask(&securityCheck)
			}
			if this.SecurityCheckParams.LeakScan {
				securityCheck := models.SecurityCheck{
					LeakScan:  this.SecurityCheckParams.LeakScan,
					Container: object,
					Type:      models.Sc_Type_Container,
					Job:       this.Job,
				}
				this.PrePareTask(&securityCheck)
			}
		}
	case models.Sc_Type_Image:
		for _, image := range this.ImageList {
			if this.SecurityCheckParams.HostImageVulnScan {
				securityCheck := models.SecurityCheck{
					HostImageVulnScan: this.SecurityCheckParams.HostImageVulnScan,
					Image:             image,
					Type:              models.Sc_Type_Image,
					Job:               this.Job,
				}
				this.PrePareTask(&securityCheck)
			}
		}
	}

	logs.Info("PrePare task end ......")
	this.GetCurrentBatchTask()
	logs.Info("################ PrePare work <<<end>>> ################")
}

func (this *SecurityScanService) PrePareTask(securityCheck *models.SecurityCheck) {
	//默认Job
	Job_Type_BM_Docker := this.DefaultJob[models.TMP_Type_BM_Docker]
	Job_Type_BM_K8S := this.DefaultJob[models.TMP_Type_BM_K8S]
	Job_Type_DockerVS := this.DefaultJob[models.TMP_Type_DockerVS]
	Job_Type_HostVS := this.DefaultJob[models.TMP_Type_HostVS]
	Job_Type_HostImageVulnScan := this.DefaultJob[models.TMP_Type_HostImageVulnScan]

	// 默认模板

	// 针对资产管理里的安全检测，使用系统默认Job
	if this.IsSystem && securityCheck.Job == nil {
		// 基线
		if securityCheck.DockerCIS {
			securityCheck.Job = Job_Type_BM_Docker
		}
		if securityCheck.KubenetesCIS {
			securityCheck.Job = Job_Type_BM_K8S
		}
		// 主机杀毒
		if securityCheck.VirusScan && securityCheck.Type == models.SC_Type_Host {
			securityCheck.Job = Job_Type_HostVS
		}
		// 容器杀毒
		if securityCheck.VirusScan && securityCheck.Type == models.Sc_Type_Container {
			securityCheck.Job = Job_Type_DockerVS
		}
		// 主机镜像扫描
		if securityCheck.HostImageVulnScan {
			securityCheck.Job = Job_Type_HostImageVulnScan
		}
	}

	if securityCheck.Job == nil {
		return
	}
	this.genTask(securityCheck)
}

func (this *SecurityScanService) genTask(securityCheck *models.SecurityCheck) {
	taskpre := "系统任务-"
	if securityCheck.Job.JobLevel == models.Job_Level_User {
		taskpre = "用户任务-"
	}
	task := new(models.Task)
	uid, _ := uuid.NewV4()
	task.Id = uid.String()

	if securityCheck.VirusScan {
		// 杀毒
		if this.SecurityCheckParams.Type == models.Sc_Type_Container {
			// 容器杀毒
			logs.Info("PrePare task, Type:  %s , Task Id: %s ......", models.TMP_Type_DockerVS, uid)
			task.SystemTemplate = this.DefaultTMP[models.TMP_Type_DockerVS]
			task.Description = taskpre + models.TMP_Type_DockerVS
			task.Container = securityCheck.Container
		} else if this.SecurityCheckParams.Type == models.SC_Type_Host {
			// 主机病毒
			task.Host = securityCheck.Host
			logs.Info("PrePare task, Type:  %s , Task Id: %s ......", models.TMP_Type_HostVS, uid)
			task.SystemTemplate = this.DefaultTMP[models.TMP_Type_HostVS]
			task.Description = taskpre + models.TMP_Type_HostVS
		}
		task.Container = securityCheck.Container
	} else if securityCheck.DockerCIS {
		//基线-Docker
		logs.Info("PrePare task, Type:  %s , Task Id: %s ......", models.TMP_Type_BM_Docker, uid)
		task.Description = taskpre + models.TMP_Type_BM_Docker
		task.SystemTemplate = this.DefaultTMP[models.TMP_Type_BM_Docker]
		task.Host = securityCheck.Host
	} else if securityCheck.KubenetesCIS {
		//基线-K8S
		logs.Info("PrePare task, Type:  %s , Task Id: %s ......", models.TMP_Type_BM_K8S, uid)
		task.Description = taskpre + models.TMP_Type_BM_K8S
		task.SystemTemplate = this.DefaultTMP[models.TMP_Type_BM_K8S]
		task.Host = securityCheck.Host
	} else if securityCheck.HostImageVulnScan {
		// 主机镜像扫描
		logs.Info("PrePare task, Type:  %s , Task Id: %s ......", models.TMP_Type_HostImageVulnScan, uid)
		task.Description = taskpre + models.TMP_Type_HostImageVulnScan
		task.SystemTemplate = this.DefaultTMP[models.TMP_Type_HostImageVulnScan]
		task.Image = securityCheck.Image
	}

	if securityCheck.Job.Type == "" {
		task.Type = models.Job_Type_Once
	} else {
		task.Type = securityCheck.Job.Type
	}

	// 集群ID
	if this.ClusterCheckObject != nil {
		task.ClusterId = this.ClusterCheckObject.Id
	}
	task.Name = taskpre + task.Id
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

func (this *SecurityScanService) saveTaskLog(task *models.Task, securityCheck *models.SecurityCheck) {
	taskLog := models.TaskLog{}
	taskRawInfo, _ := json.Marshal(task)
	taskLog.Task = string(taskRawInfo)
	taskLog.Account = securityCheck.Job.Account
	taskLog.Level = models.Log_level_Info
	taskLog.RawLog = fmt.Sprintf("Add security check task, Id: %s, Type: %s, Batch: %v, KStatus: %s",
		task.Id, task.Type, task.Batch, task.Status)
	taskLog.Add()
}

func (this *SecurityScanService) PrePareDefaultTMP() map[string]*models.SystemTemplate {
	if this.DefaultTMP == nil {
		st := new(models.SystemTemplate)
		this.DefaultTMP = st.GetDefaultTemplate()
	}
	logs.Info("PrePare Default Template: %v", this.DefaultTMP)
	return this.DefaultTMP
}

func (this *SecurityScanService) PrePareDefaultJob() map[string]*models.Job {
	if this.DefaultTMP == nil {
		defaultJob := new(models.Job)
		this.DefaultJob = defaultJob.GetDefaultJob()
	}
	logs.Info("PrePare Default Job: %v", this.DefaultJob)
	return this.DefaultJob
}

func (this *SecurityScanService) GetCurrentBatchTask() []*models.Task {
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

func (this *SecurityScanService) DeliverTask() models.Result {
	var ResultData models.Result
	ResultData = this.Check()
	if ResultData.Code != http.StatusOK {
		return ResultData
	}
	this.PrePare()
	wsDelive := ws.WSDeliverService{
		Hub:                  models.WSHub,
		Bath:                 this.Batch,
		CurrentBatchTaskList: this.CurrentBatchTaskList,
	}
	// 下发任务
	go wsDelive.DeliverTaskToNats()
	ResultData.Code = http.StatusOK
	data := make(map[string]interface{})
	data["items"] = this.CurrentBatchTaskList
	ResultData.Data = data
	return ResultData
}

func (this *SecurityScanService) ClusterCheck() models.Result {
	result := models.Result{Code: http.StatusOK}
	this.PrePareDefaultJob()
	this.PrePareDefaultTMP()

	// 获取集群内主机
	host := models.HostConfig{ClusterId: this.ClusterCheckObject.Id, ClusterName: this.ClusterCheckObject.Name}
	msg := ""
	hostResult := host.List(0, 0)
	if hostResult.Code != http.StatusOK {
		result.Code = utils.GetHostListForClusterErr
		msg = "GetHostListForClusterErr"
		result.Message = msg
		return result
	}
	if hostResult.Data == nil {
		result.Code = utils.NotFoundHostForClusterErr
		msg = "NotFoundHostForClusterErr"
		result.Message = msg
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

func (this *SecurityScanService) Check() models.Result {
	ResultData := models.Result{Code: http.StatusOK}
	baseService := base.BaseService{
		HostIds:      this.SecurityCheckParams.HostIds,
		ImageIds:     this.ImageIds,
		ContainerIds: this.ContainerIds,
		JobId:        this.SecurityCheckParams.JobId}

	// job检查
	if !this.IsSystem {
		ResultData, this.Job = baseService.CheckJobIsExist()
		if ResultData.Code != http.StatusOK {
			return ResultData
		}
	}

	// 主机授权、主机资源检查
	ResultData, this.HostList = baseService.CheckHostLicense()
	if ResultData.Code != http.StatusOK {
		return ResultData
	}

	// 资源镜像资源
	ResultData, this.ImageList = baseService.CheckImageIsExist()
	if ResultData.Code != http.StatusOK {
		return ResultData
	}

	// 检查容器资源
	ResultData, this.ContainerList = baseService.CheckContainerIsExist()
	if ResultData.Code != http.StatusOK {
		return ResultData
	}
	return ResultData
}
