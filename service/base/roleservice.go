package base

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type RoleService struct {
	Role *models.Role
}

func (this *RoleService) CheckUserInRole() models.Result {
	result := models.Result{Code: http.StatusOK}
	userQuery := models.User{Role: this.Role}
	userList, _ := userQuery.Get()
	if userList != nil {
		result.Code = utils.ForbiddenRoleRemoveErr
		result.Message = fmt.Sprintf("Has User in Role, Remove Role failed,  code: %d ", result.Code)
		logs.Error(result.Message)
		return result
	}
	return result
}
