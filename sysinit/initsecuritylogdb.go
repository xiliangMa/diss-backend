package sysinit

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
	_ "github.com/xiliangMa/diss-backend/models/securitylog"
	"github.com/xiliangMa/diss-backend/utils"
	"os"
)

func InitSecurityLogDB() {
	driver := utils.DS_Driver_Postgres
	runMode := beego.AppConfig.String("RunMode")
	envRunMode := os.Getenv("RunMode")
	if envRunMode != "" {
		runMode = envRunMode
	}
	DSAlias := utils.DS_Security_Log
	// true: drop table 后再建表
	force, _ := beego.AppConfig.Bool("Force")

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
	// 生产环境
	if runMode == utils.Run_Mode_Prod {
		//数据库名称
		dbName = os.Getenv(utils.DS_Security_Log_DB_Name)
		//数据库连接用户名
		dbUser = os.Getenv(utils.DS_Security_Log_User)
		//数据库连接密码
		dbPwd = os.Getenv(utils.DS_Security_Log_Pwd)
		//数据库IP（域名）
		dbHost = os.Getenv(utils.DS_Security_Log_Host)
	}

	// demo mysql
	//orm.RegisterDataBase("default", "mysql", "root:abc123@tcp(127.0.0.1:3306)/uranus_local?charset=utf8")
	//DS := fmt.Sprintf("%s%s%s%s%s%s%s%s", dbUser, ":", dbPwd, "@tcp(", dbHost, ":"+port+")/", dbName, "?charset=utf8")

	// postgres
	DS := fmt.Sprintf("%s%s%s%s%s%s", "host="+dbHost, " port="+port, " user="+dbUser, " password="+dbPwd, " dbname="+dbName, " sslmode=disable")
	orm.RegisterDriver(driver, orm.DRPostgres)
	err := orm.RegisterDataBase(DSAlias, driver, DS)
	if err != nil {
		logs.Error("DB Register fail, >>> DSAlias: %s <<<, Err: %s", DSAlias, err)
	}

	// auto create db
	err = orm.RunSyncdb(DSAlias, force, true)
	if err != nil {
		logs.Error("Auth Create table fail, >>> DSAlias: %s <<<, Err: %s", DSAlias, err)
	}

}
