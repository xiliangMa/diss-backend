package job

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/securitycheck"
	"github.com/xiliangMa/diss-backend/utils"
	"time"
)

type JobService struct {
	JobParam          *models.Job
	securityCheckList []*models.SecurityCheck
}

func (this *JobService) ActiveJob() models.Result {
	var ResultData models.Result
	job := this.JobParam.Get()

	if job != nil && job.Id != "" {
		// job 处于禁用状态，不执行生成，返回错误信息
		if job.Status == models.Job_Status_Disable {
			ResultData.Code = utils.JobIsDisabledErr
			msg := fmt.Sprintf("Id:%s, Job:%s is Disabled", job.Id, job.Name)
			logs.Error(msg)
			ResultData.Message = msg
			return ResultData
		}
		// 通过job内容 构建checklist
		this.securityCheckList = this.BuildCheckList(job)

		if len(this.securityCheckList) > 0 {
			// 通过checklist生成task
			securityCheckService := this.GenTask()

			// 下发task，更新job状态
			ResultData = securityCheckService.DeliverTask(false)
			job.Status = models.Job_Status_Active
			job.Update()
		}
	}

	return ResultData
}

func (this *JobService) GenTask() securitycheck.SecurityCheckService {

	batch := time.Now().Unix()
	SecCheckListModel := models.SecurityCheckList{CheckList: this.securityCheckList}
	securityCheckService := securitycheck.SecurityCheckService{SecurityCheckList: &SecCheckListModel, Batch: batch}
	securityCheckService.PrePare()
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
		seccheck.DockerBenchMarkCheck = true
	case models.TMP_Type_DockerVS, models.TMP_Type_HostVS:
		seccheck.VirusScan = true
	case models.TMP_Type_LS:
		seccheck.LeakScan = true
	}
	return seccheck
}
