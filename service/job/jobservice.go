package job

import (
	"fmt"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/securitycheck"
	"github.com/xiliangMa/diss-backend/utils"
	"time"
)

type JobService struct {
	JobParam *models.Job
}

func (this *JobService) GenTaskList(account string) models.Result {
	var ResultData models.Result

	if account == "" {
		account = models.Account_Admin
	}
	batch := time.Now().Unix()
	joblist := this.JobParam.Get(this.JobParam.Id)
	secCheckList := []*models.SecurityCheck{}
	if len(joblist) > 0 {
		job := joblist[0]

		if job.Status == models.Job_Status_Disable {
			ResultData.Code = utils.JobDisabledErr
			ResultData.Message = fmt.Sprintf("Id: %s , Name: %s Job is Disabled", job.Id, job.Name)
			return ResultData
		}
		// 通过job内容 构建checklist
		targetType := job.SystemTemplate.Type //目标类型为检查模板的类型
		resType := models.SC_Type_Host        // 资源类型包括 主机 容器 镜像
		if (targetType == models.TMP_Type_IDS_Docker) || (targetType == models.TMP_Type_DockerVS) {
			resType = models.Sc_Type_Container
		}
		if (resType == models.SC_Type_Host) && (job.HostConfig != nil) {
			for _, checkHost := range job.HostConfig {
				seccheck := new(models.SecurityCheck)
				seccheck.Type = resType
				switch targetType {
				case models.TMP_Type_BM_Docker, models.TMP_Type_BM_K8S:
					seccheck.BenchMarkCheck = true
				case models.TMP_Type_DockerVS, models.TMP_Type_HostVS:
					seccheck.VirusScan = true
				case models.TMP_Type_LS:
					seccheck.LeakScan = true
				}
				seccheck.Host = checkHost
				seccheck.Job = job
				secCheckList = append(secCheckList, seccheck)
			}
		} else if (resType == models.Sc_Type_Container) && (job.ContainerConfig != nil) {
			for _, checkContainer := range job.ContainerConfig {
				seccheck := new(models.SecurityCheck)
				seccheck.Type = resType
				switch targetType {
				case models.TMP_Type_BM_Docker, models.TMP_Type_BM_K8S:
					seccheck.BenchMarkCheck = true
				case models.TMP_Type_DockerVS, models.TMP_Type_HostVS:
					seccheck.VirusScan = true
				case models.TMP_Type_LS:
					seccheck.LeakScan = true
				}
				seccheck.Container = checkContainer
				seccheck.Job = job
				secCheckList = append(secCheckList, seccheck)
			}
		}
	}

	SecCheckListModel := models.SecurityCheckList{CheckList: secCheckList}
	securityCheckService := securitycheck.SecurityCheckService{SecurityCheckList: &SecCheckListModel, Batch: batch, Account: account}
	result := securityCheckService.DeliverTask()
	return result
}
