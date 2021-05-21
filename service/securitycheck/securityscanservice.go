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
	ClusterList          []*models.Cluster
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
					PathList:  this.SecurityCheckParams.PathList,
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
					PathList:  this.SecurityCheckParams.PathList,
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
			if this.SecurityCheckParams.ImageVulnScan {
				securityCheck := models.SecurityCheck{
					ImageVulnScan: this.SecurityCheckParams.ImageVulnScan,
					Image:         image,
					Type:          models.Sc_Type_Image,
					Job:           this.Job,
				}
				this.PrePareTask(&securityCheck)
			}
			if this.SecurityCheckParams.VirusScan {
				securityCheck := models.SecurityCheck{
					VirusScan: this.SecurityCheckParams.VirusScan,
					Image:     image,
					Type:      models.Sc_Type_Image,
					Job:       this.Job,
					PathList:  this.SecurityCheckParams.PathList,
				}
				this.PrePareTask(&securityCheck)
			}
		}
	case models.Sc_Type_Cluster:
		for _, cluster := range this.ClusterList {
			if this.SecurityCheckParams.KubenetesScan {
				securityCheck := models.SecurityCheck{
					KubenetesScan: this.SecurityCheckParams.KubenetesScan,
					Cluster:       cluster,
					Type:          models.Sc_Type_Cluster,
					Job:           this.Job,
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
	Job_Type_ContainerVS := this.DefaultJob[models.TMP_Type_ContainerVS]
	Job_Type_HostVS := this.DefaultJob[models.TMP_Type_HostVS]
	Job_Type_ImageVS := this.DefaultJob[models.TMP_Type_ImageVS]
	Job_Type_HostImageVulnScan := this.DefaultJob[models.TMP_Type_HostImageVulnScan]
	Job_Type_ImageScan := this.DefaultJob[models.TMP_Type_ImageVulnScan]
	Job_Type_KubeScan := this.DefaultJob[models.TMP_Type_KubernetesVulnScan]

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
			securityCheck.Job = Job_Type_ContainerVS
		}
		// 镜像杀毒
		if securityCheck.VirusScan && securityCheck.Type == models.Sc_Type_Image {
			securityCheck.Job = Job_Type_ImageVS
		}
		// 主机镜像扫描
		if securityCheck.HostImageVulnScan {
			securityCheck.Job = Job_Type_HostImageVulnScan
		}
		// 仓库镜像扫描
		if securityCheck.ImageVulnScan {
			securityCheck.Job = Job_Type_ImageScan
		}
		// 集群漏扫
		if securityCheck.KubenetesScan {
			securityCheck.Job = Job_Type_KubeScan
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
			logs.Info("PrePare task, Type:  %s , Task Id: %s ......", models.TMP_Type_ContainerVS, uid)
			task.SystemTemplate = this.DefaultTMP[models.TMP_Type_ContainerVS]
			task.Description = taskpre + models.TMP_Type_ContainerVS
			task.Container = securityCheck.Container
		} else if this.SecurityCheckParams.Type == models.SC_Type_Host {
			// 主机病毒
			task.Host = securityCheck.Host
			logs.Info("PrePare task, Type:  %s , Task Id: %s ......", models.TMP_Type_HostVS, uid)
			task.SystemTemplate = this.DefaultTMP[models.TMP_Type_HostVS]
			task.Description = taskpre + models.TMP_Type_HostVS
		} else if this.SecurityCheckParams.Type == models.Sc_Type_Image {
			// 主机病毒
			task.Image = securityCheck.Image
			logs.Info("PrePare task, Type:  %s , Task Id: %s ......", models.TMP_Type_ImageVS, uid)
			task.SystemTemplate = this.DefaultTMP[models.TMP_Type_ImageVS]
			task.Description = taskpre + models.TMP_Type_ImageVS
		}

		if task.SystemTemplate.DefaultPathList != "" {
			task.PathList = task.Job.SystemTemplate.DefaultPathList
		}
		if securityCheck.PathList != "" {
			task.PathList = securityCheck.PathList
		}
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
	} else if securityCheck.ImageVulnScan {
		// 仓库镜像扫描
		logs.Info("PrePare task, Type:  %s , Task Id: %s ......", models.TMP_Type_ImageVulnScan, uid)
		task.Description = taskpre + models.TMP_Type_ImageVulnScan
		task.SystemTemplate = this.DefaultTMP[models.TMP_Type_ImageVulnScan]
		task.Image = securityCheck.Image
	} else if securityCheck.KubenetesScan {
		//kubernetes 漏扫
		logs.Info("PrePare task, Type:  %s , Task Id: %s ......", models.TMP_Type_KubernetesVulnScan, uid)
		task.Description = taskpre + models.TMP_Type_KubernetesVulnScan
		task.SystemTemplate = this.DefaultTMP[models.TMP_Type_KubernetesVulnScan]
		task.ClusterOBJ = securityCheck.Cluster
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
	task.Status = models.Task_Status_Created
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
	taskLog.RawLog = fmt.Sprintf("Add task success, Status: %s, Type: %s, Batch: %v, TaskId: %s.", task.Status, task.Type, task.Batch, task.Id)
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
		JobId:        this.SecurityCheckParams.JobId,
		ClusterIds:   this.SecurityCheckParams.ClusterIds,
	}
	// 1. 资源检查
	// 镜像资源
	ResultData, this.ImageList = baseService.CheckImageIsExist()
	if ResultData.Code != http.StatusOK {
		return ResultData
	}

	// 容器资源
	ResultData, this.ContainerList = baseService.CheckContainerIsExist()
	if ResultData.Code != http.StatusOK {
		return ResultData
	}

	// kubernetes 资源
	ResultData, this.ClusterList, baseService.HostListInCluster = baseService.CheckClusterIsExist()
	if ResultData.Code != http.StatusOK {
		return ResultData
	}

	// 2. job检查
	if !this.IsSystem {
		ResultData, this.Job = baseService.CheckJobIsExist()
		if ResultData.Code != http.StatusOK {
			return ResultData
		}
	}

	// 3. 主机授权检查
	if !this.ImageVulnScan && !this.KubenetesScan {
		ResultData, this.HostList = baseService.CheckHostLicense()
		if ResultData.Code != http.StatusOK {
			return ResultData
		}
	}

	return ResultData
}
