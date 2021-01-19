package models

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/xiliangMa/diss-backend/utils"
)

// 用户-租户绑定关系 diss_api
type AccountUsers struct {
	UserName    string `description:"(用户名)"`
	AccountName string `description:"(租户)"`
	Type        string `description:"(类型)"`
	Source      string `description:"(资源)"`
	CreatedAt   int    `description:"(创建时间)"`
	LastUpdated int    `description:"(更新时间)"`
	uuid        string `description:"(id)"`
}

type AccountUsersInterface interface {
	GetAccountByUser()
}

func (this *AccountUsers) GetAccountByUser() (error, string) {
	o := orm.NewOrmUsingDB(utils.DS_Diss_Api)
	var accountName string

	err := o.Raw("select account_name from "+utils.AccountUsers+" where username = ?", this.UserName).QueryRow(&accountName)
	if err != nil {
		if err.Error() == "<QuerySeter> no row found" {
			logs.Error("Get AccountUsers List failed, code: %d, err: %s", utils.NoAccountUsersErr, err)
		} else {
			logs.Error("Get AccountUsers List failed, code: %d, err: %s", utils.GetAccountUsersErr, err)
		}

	}
	return err, accountName
}
