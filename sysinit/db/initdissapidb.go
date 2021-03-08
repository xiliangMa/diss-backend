package db

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
	_ "github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
)

type DissApiDB struct {
}

func (this *DissApiDB) InitDB() {
	driver := utils.DS_Driver_Postgres
	DSAlias := utils.DS_Diss_Api
	// true: drop table 后再建表
	//force, _ := beego.AppConfig.Bool("Force")

	//连接名称
	//dbAlias := beego.AppConfig.String(DS + "::Alias")
	//数据库名称
	dbName := beego.AppConfig.String(DSAlias + "::DBName")
	//数据库连接用户名
	dbUser := beego.AppConfig.String(DSAlias + "::User")
	//数据库连接密码
	dbPwd := beego.AppConfig.String(DSAlias + "::Pwd")
	//数据库IP（域名）
	dbHost := beego.AppConfig.String(DSAlias + "::Host")
	//端口
	port := beego.AppConfig.String(DSAlias + "::Port")
	// postgres
	DS := fmt.Sprintf("%s%s%s%s%s%s", "host="+dbHost, " port="+port, " user="+dbUser, " password="+dbPwd, " dbname="+dbName, " sslmode=disable")
	orm.RegisterDriver(driver, orm.DRPostgres)
	err := orm.RegisterDataBase(DSAlias, driver, DS)
	orm.SetMaxIdleConns(DSAlias, 2)
	orm.SetMaxOpenConns(DSAlias, 100)
	if err != nil {
		logs.Error("DB Register fail, >>> DSAlias: %s <<<, Err: %s", DSAlias, err)
	}

	//auto create db
	//err = orm.RunSyncdb(DSAlias, force, true)
	//if err != nil {
	//	logs.Error("Auth Create table fail, >>> DSAlias: %s <<<, Err: %s", DSAlias, err)
	//}
	//// 设置为 UTC 时间
	//orm.DefaultTimeLoc = time.UTC
}
