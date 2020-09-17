package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/auth"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type AuthController struct {
	beego.Controller
}

// @Title login
// @Description login
// @Param name query string true "userName"
// @Param pwd query string true "userPwd"
// @Success 200 {object} models.Result
// @router /login [post]
func (this *AuthController) Login() {
	name := this.GetString("name")
	pwd := this.GetString("pwd")

	var ResultData models.Result
	if name == "" || pwd == "" {
		ResultData.Message = "UserAndPwdNotNull"
		ResultData.Code = utils.UserAndPwdNotNull
		this.Data["json"] = ResultData
		this.ServeJSON(false)
		return

	}
	jwtService := auth.JwtService{}
	result, code := jwtService.GreateToken(name, pwd)
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
