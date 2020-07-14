package db

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/sysinit/dbscript"
	"github.com/xiliangMa/diss-backend/utils"
	"os"
)

type DefaultDB struct {
}

func (this *DefaultDB) InitDB() {
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

	this.registerModel()
	//auto create db
	err = orm.RunSyncdb(DSAlias, force, true)
	if err != nil {
		logs.Error("Auth Create table fail, >>> DSAlias: %s <<<, Err: %s", DSAlias, err)
	}
	// 设置为 UTC 时间(默认为Local时间 需要可以更改)
	//orm.DefaultTimeLoc = time.UTC

	//初始化系统数据
	this.InitSystemData()
}

func (this *DefaultDB) registerModel() {
	// k8s
	orm.RegisterModel(new(models.Cluster), new(models.NameSpace), new(models.Pod), new(models.Service), new(models.Deployment))
	// securitylog
	orm.RegisterModel(new(models.BenchMarkLog))
	// task
	orm.RegisterModel(new(models.Task), new(models.Job))
	//base
	orm.RegisterModel(new(models.SystemTemplate), new(models.SystemTemplateGroup),
		new(models.ContainerConfig), new(models.ContainerInfo), new(models.ContainerPs),
		new(models.HostConfig), new(models.HostInfo), new(models.HostPs), new(models.ImageConfig),
		new(models.ImageInfo), new(models.Groups))
	// logconfig and timeEdgepoint
	orm.RegisterModel(new(models.LogConfig), new(models.TimeEdgePoint))
	// license
	orm.RegisterModel(new(models.LicenseConfig), new(models.LicenseModule), new(models.LicenseHistory))
}

func (this *DefaultDB) InitSystemData() {
	// 初始化
	o := orm.NewOrm()
	logs.Info("Init default data start................")

	logs.Info("Init default SystemTemplate >>>>>>>>>>>>>>>>")
	// 系统魔板
	_, err := o.Raw(dbscript.DefaultDockerBenchSql).Exec()
	if err != nil {
		logs.Error("Init DefaultDockerBench fail, err: %s", err)
	}

	_, err = o.Raw(dbscript.DefaultK8sBenchSql).Exec()
	if err != nil {
		logs.Error("Init DefaultK8sBench fail, err: %s", err)
	}

	_, err = o.Raw(dbscript.DefaultDockerVirusScanSql).Exec()
	if err != nil {
		logs.Error("Init DefaultDockerVirusScan fail, err: %s", err)
	}

	_, err = o.Raw(dbscript.DefaultHostrVirusScanSql).Exec()
	if err != nil {
		logs.Error("Init DefaultHostrVirusScanSql fail, err: %s", err)
	}

	logs.Info("Init default system Job >>>>>>>>>>>>>>>>")

	//系统job
	_, err = o.Raw(dbscript.DefaultDockerBenchJobSql).Exec()
	if err != nil {
		logs.Error("Init DefaultDockerBench Job fail, err: %s", err)
	}

	_, err = o.Raw(dbscript.DefaultK8sBenchJobSql).Exec()
	if err != nil {
		logs.Error("Init DefaultK8sBench Job fail, err: %s", err)
	}

	_, err = o.Raw(dbscript.DefaultDockerVirusScanJobSql).Exec()
	if err != nil {
		logs.Error("Init DefaultDockerVirusScan Job fail, err: %s", err)
	}

	_, err = o.Raw(dbscript.DefaultHostVirusScanJobSql).Exec()
	if err != nil {
		logs.Error("Init DefaultHostVirusScanJobSql Job fail, err: %s", err)
	}

	logs.Info("Init default data end................")
}
