package base

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type UserService struct {
	UserId    int
	LoginName string
	IsAdmin   bool
}

func (this *UserService) ConstrictUserModify() models.Result {
	result := models.Result{Code: http.StatusOK}
	userQuery := models.User{Id: this.UserId}
	userObject, _ := userQuery.Get()

	if userObject.Name != this.LoginName {
		result.Code = utils.ForbidenModifyUserErr
		result.Message = fmt.Sprintf("Forbinden Opeartion: Modify User failed,  code: %d ", result.Code)
		logs.Error(result.Message)
		return result
	}
	return result
}

func (this *UserService) ConstrictUserRemove() models.Result {
	result := models.Result{Code: http.StatusOK}
	if this.IsAdmin {
		userQuery := models.User{Id: this.UserId}
		userObject, _ := userQuery.Get()
		if userObject.Name == this.LoginName {
			result.Code = utils.ForbidenAdminDelSelfErr
			result.Message = fmt.Sprintf("Forbinden Opeartion: Admin Remove Self failed,  code: %d ", result.Code)
		}
		return result
	}

	result.Code = utils.ForbidenRemoveUserErr
	result.Message = fmt.Sprintf("Forbinden Opeartion: Remove User failed,  code: %d ", result.Code)
	logs.Error(result.Message)

	return result
}

func (this *UserService) ConstrictUserCreate() models.Result {
	result := models.Result{Code: http.StatusOK}
	result.Code = utils.ForbidenCreateUserErr
	result.Message = fmt.Sprintf("Forbinden Opeartion: Create User failed,  code: %d ", result.Code)
	logs.Error(result.Message)

	return result
}
