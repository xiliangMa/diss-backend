package models

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/xiliangMa/diss-backend/utils"
	"strings"
)

type Registries struct {
	Registry     string `description:"(仓库)"`
	RegistryType string `description:"(仓库类型)"`
	UserId       string `description:"(租户)"`
}

type RegistriesInterface interface {
	List(from, limit int) []*Registries
}

func (this *Registries) List() []*Registries {
	o := orm.NewOrmUsingDB(utils.DS_Diss_Api)
	var registList []*Registries = nil
	var ResultData Result
	var err error

	sql := ` select * from registries `
	filter := ""
	if this.Registry != "" {
		filter = filter + `registry = '` + this.Registry + `' and `
	}
	if this.RegistryType != "" {
		filter = filter + `registry_type  = '` + this.RegistryType + `' and `
	}
	if this.UserId != "" {
		filter = filter + `registries."userId" like '%` + this.UserId + `%' and `
	}

	if filter != "" {
		sql = sql + " where " + filter
	}
	sql = strings.TrimSuffix(strings.TrimSpace(sql), "and")

	_, err = o.Raw(sql).QueryRows(&registList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetRegistriesErr
		logs.Error("Get Registries List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return nil
	}

	return registList
}
