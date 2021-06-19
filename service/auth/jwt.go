package auth

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
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
	if name == beego.AppConfig.String("system::AdminUser") && pwd == beego.AppConfig.String("system::AdminPwd") {
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
		return TokenStr, http.StatusOK
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
