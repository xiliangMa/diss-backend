package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/auth"
	"github.com/xiliangMa/diss-backend/service/system/system"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type AuthController struct {
	web.Controller
}

// @Title login
// @Description login
// @Param name query string true "userName"
// @Param pwd query string true "userPwd"
// @Param loginType query string true "loginType类型 LOCAL(本地) LDAP(ldap登录)，空值视为LOCAL"
// @Success 200 {object} models.Result
// @router /login [post]
func (this *AuthController) Login() {
	name := this.GetString("name")
	pwd := this.GetString("pwd")
	loginType := this.GetString("loginType")

	var ResultData models.Result
	if name == "" || pwd == "" {
		ResultData.Message = "UserAndPwdNotNull"
		ResultData.Code = utils.UserAndPwdNotNull
		this.Data["json"] = ResultData
		this.ServeJSON(false)
		return

	}
	jwtService := auth.JwtService{}
	jwtService.LoginType = loginType
	result, code := jwtService.CreateToken(name, pwd)
	ResultData.Code = code
	if code != http.StatusOK {
		ResultData.Message = result
	} else {
		ResultData.Data = result
	}
	this.Data["json"] = ResultData
	this.ServeJSON(false)
}

// @Title authorization
// @Description authorization
// @Param token header string true "authToken"
// @Success 200 {object} models.Result
// @router /Authorization [post]
func (this *AuthController) Authorize() {

	token := this.Ctx.Request.Header.Get("token")
	var ResultData models.Result
	jwtService := auth.JwtService{}
	result, code := jwtService.CheckToken(token)
	ResultData.Code = code
	if code != http.StatusOK {
		ResultData.Message = result
		logs.Error("Authorization failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
	} else {
		logs.Info("login success, %s", ResultData)
		ResultData.Data = result
	}
	this.Data["json"] = ResultData
	this.ServeJSON(false)
}

// @Title CheckLDAPStatus
// @Description Check LDAP Status
// @Success 200 {object} models.Result
// @router /checkldap [post]
func (this *AuthController) CheckLDAPStatus() {
	ldapService := system.LDAPService{}
	this.Data["json"] = ldapService.CheckLDAPStatus()
	this.ServeJSON(false)
}
