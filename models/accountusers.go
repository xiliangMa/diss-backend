package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"time"
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

func (this *AccountUsers) GetAccountByUser() (error, *AccountUsers) {
	o := orm.NewOrm()
	orm.DefaultTimeLoc = time.Local
	o.Using(utils.DS_Diss_Api)
	var accountUsers AccountUsers

	err := o.Raw("select * from "+utils.AccountUsers+" user_name = ?", this.UserName).QueryRow(&accountUsers)
	if err != nil {
		logs.Error("Get AccountUsers List failed, code: %d, err: %s", utils.GetAccountUsersErr, "GetAccountUsersErr")
	}
	return err, &accountUsers
}
