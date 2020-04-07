package db

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
	"github.com/xiliangMa/diss-backend/models"
	mk8s "github.com/xiliangMa/diss-backend/models/k8s"
	msecuritylog "github.com/xiliangMa/diss-backend/models/securitylog"
	mtask "github.com/xiliangMa/diss-backend/models/task"
	"github.com/xiliangMa/diss-backend/utils"
	"os"
)

func InitDB() {
	driver := utils.DS_Driver_Postgres
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
		dbName = os.Getenv(utils.DS_Default_POSTGRES_DB)
		//数据库连接用户名
		dbUser = os.Getenv(utils.DS_Default_POSTGRES_USER)
		//数据库连接密码
		dbPwd = os.Getenv(utils.DS_Default_POSTGRES_PASSWORD)
		//数据库IP（域名）
		dbHost = os.Getenv(utils.DS_Default_POSTGRES_HOST)
		//端口
		port = os.Getenv(utils.DS_Default_POSTGRES_PORT)
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

	registerDefaultModel()
	//auto create db
	err = orm.RunSyncdb(DSAlias, force, true)
	if err != nil {
		logs.Error("Auth Create table fail, >>> DSAlias: %s <<<, Err: %s", DSAlias, err)
	}
	// 设置为 UTC 时间(默认为Local时间 需要可以更改)
	//orm.DefaultTimeLoc = time.UTC
}

func registerDefaultModel() {
	// k8s
	orm.RegisterModel(new(mk8s.Cluster), new(mk8s.NameSpace), new(mk8s.Pod))
	// securitylog
	orm.RegisterModel(new(msecuritylog.BenchMarkLog))
	// task
	orm.RegisterModel(new(mtask.Task))
	//base
	orm.RegisterModel(new(models.BenchMarkTemplate), new(models.CmdHistory),
		new(models.ContainerConfig), new(models.ContainerInfo), new(models.ContainerPs),
		new(models.HostConfig), new(models.HostInfo), new(models.HostPs), new(models.ImageConfig),
		new(models.ImageInfo), new(models.Groups))
}
