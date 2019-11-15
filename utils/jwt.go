package utils

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

var (
	secret []byte = []byte("secret")
)

func GreateToken(name, pwd string) (string, int) {
	if name == beego.AppConfig.String("system::AdminUser") && pwd == beego.AppConfig.String("system::AdminPwd") {
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = "Jon Snow"
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err.Error(), SiginErr
		}
		return t, http.StatusOK
	}
	return http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized
}

func CheckToken(tokenString string) (string, int) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return "", http.StatusOK
	} else {
		return err.Error(), AuthorizeErr
	}
}
