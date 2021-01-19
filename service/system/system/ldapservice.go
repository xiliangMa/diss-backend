package system

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/go-ldap/ldap"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type LDAPService struct {
	LoginUser models.LDAPUser
}

func (this *LDAPService) UserAuthentication() (bool, error) {

	lc := &models.LM.Conn
	if *lc == nil {
		// 自动执行重连一次
		models.LM.CreateLDAPConnection()
		if *lc == nil {
			// 如果重连接后依然失败，返回错误
			return false, models.LM.ConnectErr
		}
	}

	logined := false
	sr, err := this.SearchUser()

	if sr != nil && len(sr.Entries) > 0 {
		// 如果没有数据返回或者超过1条数据返回，这对于用户认证而言都是不允许的。
		// 前者意味着没有查到用户，后者意味着存在重复数据
		if len(sr.Entries) > 1 {
			logs.Warn("Get More than One User.")
		}
		userdn := sr.Entries[0].DN
		password := this.LoginUser.Password

		err = (*lc).Bind(userdn, password)
		if err != nil {
			logs.Error("Bind login user fail, Err: %s.", err)
			return logined, err
		}
		logined = true
	}

	return logined, err
}

func (this *LDAPService) CheckLDAPStatus() models.Result {
	ResultData := models.Result{Code: http.StatusOK}
	ldapStatus := true
	if (!models.LM.Config.Enable) || (models.LM.Conn == nil) || (models.LM.ConnectErr != nil) {
		ldapStatus = false
		ResultData.Code = utils.LDAPConnErr
		// LDAP 服务连接错误的详细内容
		if models.LM.ConnectErr != nil {
			ResultData.Message = models.LM.ConnectErr.Error()
		}
	}
	ResultData.Data = ldapStatus
	return ResultData
}

func (this *LDAPService) SearchUser() (*ldap.SearchResult, error) {
	lcon := models.LM.Conn
	if lcon == nil {
		return nil, models.LM.ConnectErr
	}

	searchRequest := ldap.NewSearchRequest(
		models.LM.Config.BaseDn,
		// 参数分别是 scope, derefAliases, sizeLimit, timeLimit,  typesOnly [RFC4511]
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(models.LM.Config.AuthFilter, this.LoginUser.UserName),
		[]string{"dn"},
		nil,
	)

	sr, err := lcon.Search(searchRequest)
	if err != nil {
		logs.Error("Search User fail, error: %s.", err)
	}
	return sr, err
}
