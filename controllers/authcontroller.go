package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/auth"
	"github.com/xiliangMa/diss-backend/service/system/system"
	"net/http"
)

type AuthController struct {
	beego.Controller
}

// @Title login
// @Description login
// @Param name query string true "userName"
// @Param pwd query string true "userPwd"
// @Param loginType query string true "loginType类型 LOCAL(本地) LDAP(ldap登录)，空值视为LOCAL"
// @Success 200 {object} models.Result
// @router /login [post]
func (this *AuthController) Login() {
	ResultData := models.Result{Code: http.StatusOK}
	name := this.GetString("name")
	pwd := this.GetString("pwd")
	loginType := this.GetString("loginType")
	jwtService := auth.JwtService{}
	result, code := jwtService.CreateToken(name, pwd, loginType)
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
