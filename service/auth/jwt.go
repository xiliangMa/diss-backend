package auth

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/system/system"
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
}

func (this *JwtService) CreateToken(name, pwd, userType string) (string, int) {
	//检测 diss-api 用户
	loginUser := models.UserAccessCredentials{UserName: name}
	user := &models.UserAccessCredentials{}

	if this.LoginType != models.Login_Type_LDAP {
		if beego.AppConfig.String("RunMode") == "prod" {
			user = loginUser.Get()
			if user == nil {
				return "User Not Found", utils.NoSuchUser
			}
			if userType == models.Login_Type_LOCAL {
				match, err := utils.ComparePassword(pwd, user.Value)
				if !match || err != nil {
					return "Password Invalid", utils.SiginErr
				}
			} else if userType != models.Login_Type_DEV {
				return "Login Type Error", utils.LoginTypeErr
			}
		} else {
			if userType != models.Login_Type_LOCAL && userType != models.Login_Type_DEV {
				return "Login Type Error", utils.LoginTypeErr
			}
		}
	} else {
		if !models.LM.Config.Enable {
			return "LDAP is Disabled", utils.LDAPIsDisabledErr
		}
		ldapService := system.LDAPService{}
		ldapService.LoginUser = models.LDAPUser{
			UserName: name,
			Password: pwd,
		}
		logined, err := ldapService.UserAuthentication()
		if err != nil {
			return err.Error(), utils.LDAPLoginErr
		}
		if logined {
			user = &models.UserAccessCredentials{
				UserName: ldapService.LoginUser.UserName,
				Type:     models.Login_Type_LDAP,
			}
		} else {
			user = nil
		}
	}
	if user != nil || (name == beego.AppConfig.String("system::AdminUser") && pwd == beego.AppConfig.String("system::AdminPwd")) {
		// Create token
		this.Token = jwt.New(jwt.SigningMethodHS256)

		claims := this.Token.Claims.(jwt.MapClaims)
		claims["name"] = beego.AppConfig.String("AppName")
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
		claims["UserName"] = name
		claims["Pwd"] = pwd

		t, err := this.Token.SignedString([]byte("secret"))
		if err != nil {
			return err.Error(), utils.SiginErr
		}
		return t, http.StatusOK
	}
	return http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized
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
		return "", http.StatusOK
	} else {
		this.Token = nil
		return "AuthorizeErr", utils.AuthorizeErr
	}
}

func (this *JwtService) GetUserFromToken() string {
	this.CheckToken(this.TokenStr)
	user := ""
	if this.Token != nil {
		user = this.Token.Claims.(jwt.MapClaims)["UserName"].(string)
	}
	return user
}
