package auth

import (
	"encoding/base64"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/dgrijalva/jwt-go"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/plugins"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

var (
	secret []byte = []byte("secret")
)

type JwtService struct {
	Token     *jwt.Token
	TokenStr  string
	LoginType string
	User      models.User
}

func (this *JwtService) CreateToken(name, pwd, userType string) models.Result {
	result := models.Result{Code: http.StatusOK}
	data := make(map[string]interface{})
	var user models.User
	user.Name = name
	password := utils.MD5(pwd)
	passwordBase64 := base64.StdEncoding.EncodeToString([]byte(password))
	user.Password = passwordBase64
	userList, count, _ := user.UserList(0, 1)

	if count > 0 {
		// Create token
		this.Token = jwt.New(jwt.SigningMethodHS256)

		claims := this.Token.Claims.(jwt.MapClaims)
		claims["name"] = beego.AppConfig.String("AppName")
		claims["admin"] = false
		claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
		claims["UserName"] = name
		claims["Pwd"] = pwd

		userRoleList, _ := models.GlobalCasbin.Enforcer.GetRolesForUser(name)
		for _, userRole := range userRoleList {
			if userRole == utils.GetRoleString("admin") {
				claims["admin"] = true
				break
			}
		}

		t, err := this.Token.SignedString([]byte("secret"))
		if err != nil {
			result.Code = utils.SiginErr
			result.Message = err.Error()
			return result
		}
		user = *userList[0]
		user.Password = ""
		data["token"] = t
		data["user"] = user
		result.Data = data
		return result
	}

	//result.Code = http.StatusUnauthorized
	//result.Message = http.StatusText(http.StatusUnauthorized)
	result.Code = utils.UsernameOrPasswordErr
	msg := fmt.Sprintf("Username or Password error, Login failed, code: %d , name: %s.", result.Code, name)
	result.Message = msg
	logs.Info(result.Message)
	return result
}

func (this *JwtService) CheckToken(TokenStr string) (string, int) {
	var err error
	this.Token, err = jwt.Parse(TokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		this.Token = nil
		return err.Error(), utils.AuthorizeErr
	}

	if _, ok := this.Token.Claims.(jwt.MapClaims); ok && this.Token.Valid {
		return TokenStr, http.StatusOK
	} else {
		this.Token = nil
		return "AuthorizeErr", utils.AuthorizeErr
	}
}

func (this *JwtService) GetUserFromToken() (string, bool) {
	this.CheckToken(this.TokenStr)
	loginUser := ""
	isadmin := false
	if this.Token != nil {
		isadmin = this.Token.Claims.(jwt.MapClaims)["admin"].(bool)
		loginUser = this.Token.Claims.(jwt.MapClaims)["UserName"].(string)
	}
	return loginUser, isadmin
}

func (this *JwtService) CheckCaptcha() models.Result{
	var ResultData models.Result
	ResultData.Code = http.StatusOK

	// 验证码检查
	captchaVeri := new(models.Captcha)
	captchaTool := plugins.CaptchaTool{}
	captchaVeri.Id = this.User.CaptchaId
	captchaVeri.Value = this.User.CaptchaValue
	if captchaVeri.Id == "" || captchaVeri.Value == "" {
		ResultData.Message = "No Verify Code Input"
		ResultData.Code = utils.NoLoginVerifyCodeErr
		return ResultData
	}
	captchaTool.Captcha = *captchaVeri
	verifyStatus := captchaTool.VerifyCaptcha()
	if !verifyStatus {
		ResultData.Message =  "Verify Code Error."
		ResultData.Code = utils.LoginVerifyCodeErr
		return ResultData
	}
	return ResultData
}
