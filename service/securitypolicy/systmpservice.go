package securitypolicy

import (
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
)

type SystemTemplateService struct {
	SystemTemplate *models.SystemTemplate
}

func (this *SystemTemplateService) AddSystemTemplate() models.Result {
	var ResultData models.Result

	// 基线模板自动填充配置的预置json串
	if this.SystemTemplate.Type == models.TMP_Type_BM_Docker || this.SystemTemplate.Type == models.TMP_Type_BM_K8S {
		benchTemplate, err := this.SystemTemplate.Get()
		if err != nil {
			ResultData.Message = err.Error()
			ResultData.Code = utils.GetSYSTemplateErr
			return ResultData
		}
		if benchTemplate.Id != "" {
			this.SystemTemplate.CheckMasterJson = benchTemplate.CheckMasterJson
			this.SystemTemplate.CheckNodeJson = benchTemplate.CheckNodeJson
			this.SystemTemplate.CheckControlPlaneJson = benchTemplate.CheckControlPlaneJson
			this.SystemTemplate.CheckEtcdJson = benchTemplate.CheckEtcdJson
			this.SystemTemplate.CheckPoliciesJson = benchTemplate.CheckPoliciesJson
			this.SystemTemplate.CheckManagedServicesJson = benchTemplate.CheckManagedServicesJson
			this.SystemTemplate.Commands = benchTemplate.Commands
		}
	}
	ResultData = this.SystemTemplate.Add()

	return ResultData
}
