package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"strings"
)

type UserAccessCredentials struct {
	UserName string `orm:"column(username)" description:"(用户名)" json:"username"`
	Type     string `orm:"" description:"(loginType类型 LOCAL(本地) LDAP(ldap登录)，空值视为LOCAL)" json:"loginType"`
	Value    string `orm:"" description:"(密码)" json:"pwd"`
	CreateAt int64  `orm:"" description:"(创建时间)" json:"-"`
}

type UserAccessCredentialsInterface interface {
	Get() *UserAccessCredentials
}

func (this *UserAccessCredentials) Get() *UserAccessCredentials {
	o := orm.NewOrm()
	o.Using(utils.DS_Diss_Api)
	var userAccessCredentials *UserAccessCredentials
	var ResultData Result

	sql := ` select * from ` + string(utils.UserAccessCredentials) + ` `
	filter := ""
	if this.UserName != "" {
		filter = filter + `username = '` + this.UserName + `' and `
	}
	if this.Value != "" {
		filter = filter + `value = '` + this.Value + `' and `
	}

	if filter != "" {
		sql = sql + " where " + filter
	}
	sql = strings.TrimSuffix(strings.TrimSpace(sql), "and")
	resultSql := sql
	err := o.Raw(resultSql).QueryRow(&userAccessCredentials)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.NotFoundUserAccessCredentialsErr
		logs.Error("Not FoundUserAccessCredentials, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return nil
	}

	return userAccessCredentials
}
