package system

import (
	"fmt"
	"github.com/xiliangMa/diss-backend/models"
	"net/http"
)

type RuleDefineService struct {
	RuleDefine *models.RuleDefine
}

func (this *RuleDefineService) AddRule() models.Result {
	var result models.Result
	ruledefineObj := models.RuleDefine{}

	if this.RuleDefine != nil {
		result = this.RuleDefine.Add()
		if result.Code != http.StatusOK {
			return result
		}
	}

	result.Message = fmt.Sprintf("Added ruledefine success, name : %s ", ruledefineObj.RuleName)
	result.Code = http.StatusOK
	return result
}
