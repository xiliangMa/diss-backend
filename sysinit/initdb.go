package sysinit

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/xiliangMa/diss-backend/models"
	_ "github.com/xiliangMa/diss-backend/models/k8s"
	_ "github.com/xiliangMa/diss-backend/models/securitylog"
	_ "github.com/xiliangMa/diss-backend/models/task"
	"github.com/xiliangMa/diss-backend/utils"
	"os"
)

func InitDB() {

	runMode := beego.AppConfig.String("RunMode")
	envRunMode := os.Getenv("RunMode")
	if envRunMode != "" {
		runMode = envRunMode
	}
	DSAlias := utils.DS_Default
	// true: drop table 后再建表
	force, _ := beego.AppConfig.Bool("Force")

	//连接名称
	//dbAlias := beego.AppConfig.String(DS + "::Alias")
	//数据库名称
	dbName := beego.AppConfig.String(DSAlias + "::Name")
	//数据库连接用户名
	dbUser := beego.AppConfig.String(DSAlias + "::User")
	//数据库连接用户名
	dbPwd := beego.AppConfig.String(DSAlias + "::Pwd")
	//数据库IP（域名）
	dbHost := beego.AppConfig.String(DSAlias + "::Host")
	//端口
	port := beego.AppConfig.String(DSAlias + "::Port")
	// 生产环境
	if runMode == utils.Run_Mode_Prod {
		//数据库名称
		dbName = os.Getenv("MYSQL_DATABASE")
		//数据库连接用户名
		dbUser = os.Getenv("MYSQL_USER")
		//数据库连接用户名
		dbPwd = os.Getenv("MYSQL_PASSWORD")
		//数据库IP（域名）
		dbHost = os.Getenv("MYSQL_HOST")
	}

	// demo
	//orm.RegisterDataBase("default", "mysql", "root:abc123@tcp(127.0.0.1:3306)/uranus_local?charset=utf8")

	DS := fmt.Sprintf("%s%s%s%s%s%s%s%s", dbUser, ":", dbPwd, "@tcp(", dbHost, ":"+port+")/", dbName, "?charset=utf8")

	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase(DSAlias, "mysql", DS)
	if err != nil {
		logs.Error("DB Register fail, DSAlias: %s, Err: %s", DSAlias, err)
	}

	//auto create db
	orm.RunSyncdb(DSAlias, force, false)
}
