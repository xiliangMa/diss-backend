package models

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/go-ldap/ldap"
)

type LDAPManager struct {
	Conn       *ldap.Conn
	Config     *LDAPConfig
	ConnectErr error
}

type LDAPConfig struct {
	Enable     bool
	Addr       string
	Port       int
	BaseDn     string
	BindDn     string
	BindPass   string
	AuthFilter string
	Attributes []string
	TLS        bool
	StartTLS   bool
}

type LDAPUser struct {
	UserName string
	Password string
}

func NewLDAPManager() *LDAPManager {
	ldapManager := &LDAPManager{}
	ldapManager.CreateLDAPConnection()
	return ldapManager
}

func (this *LDAPManager) CreateLDAPConnection() {
	this.ConnectErr = nil
	this.Conn = nil

	sysConfig := SysConfig{}
	sysConfig.Key = LDAPClientConfig
	ldapClientConfigData := sysConfig.Get()

	if ldapClientConfigData != nil {
		ldapClientConfigStr := ldapClientConfigData.Value
		ldapClientConfig := &LDAPConfig{}
		err := json.Unmarshal([]byte(ldapClientConfigStr), ldapClientConfig)
		if err != nil {
			logs.Error("Encode ldapClientConfig json fail, error: %s.", err)
			return
		}
		this.Config = ldapClientConfig
	}
	if this.Config == nil || !this.Config.Enable {
		return
	}

	if this.Config.Addr != "" && this.Config.Port != 0 && this.Config.BaseDn != "" {
		logs.Info("LDAP Client Config: Server : %s:%d Base dn: %s", this.Config.Addr, this.Config.Port, this.Config.BaseDn)
		lcon, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", this.Config.Addr, this.Config.Port))
		if err != nil {
			this.ConnectErr = err
			logs.Error("Connect LDAP fail, error: %s.", err)
			return
		}

		err = lcon.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			this.ConnectErr = err
			logs.Error("Start LDAP TLS fail, error: %s.", err)
			return
		}

		err = lcon.Bind(this.Config.BindDn, this.Config.BindPass)
		if err != nil {
			this.ConnectErr = err
			logs.Error("Bind LDAP bindDn fail, error: %s.", err)
			return
		}
		this.Conn = lcon
	}
}
