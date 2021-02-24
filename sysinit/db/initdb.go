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
)

type DefaultDB struct {
}

func (this *DefaultDB) InitDB() {
	driver := utils.DS_Driver_Postgres
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
	orm.RegisterModel(new(models.Cluster), new(models.NameSpace), new(models.Pod), new(models.Service),
		new(models.Deployment), new(models.NetworkPolicy),
		new(models.CmdHistory), new(models.WarningInfo), new(models.DockerEvent), new(models.TaskLog), new(models.HostPackage))
	// securitylog
	orm.RegisterModel(new(models.BenchMarkLog))
	// task
	orm.RegisterModel(new(models.Task), new(models.Job))
	//base
	orm.RegisterModel(new(models.SystemTemplate), new(models.SystemTemplateGroup),
		new(models.ContainerConfig), new(models.ContainerInfo), new(models.ContainerPs),
		new(models.HostConfig), new(models.HostInfo), new(models.HostPs), new(models.ImageConfig),
		new(models.ImageInfo), new(models.Groups), new(models.SysConfig))
	orm.RegisterModel(new(models.LogConfig), new(models.TimeEdgePoint), new(models.WarningWhiteList))
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

	_, err = o.Raw(dbscript.DefaultK8sBenchCis13Sql).Exec()
	if err != nil {
		logs.Error("Init DefaultK8sBench - 1.3 fail, err: %s", err)
	}

	_, err = o.Raw(dbscript.DefaultK8sBenchCis14Sql).Exec()
	if err != nil {
		logs.Error("Init DefaultK8sBench - 1.4 fail, err: %s", err)
	}

	_, err = o.Raw(dbscript.DefaultK8sBenchCis15Sql).Exec()
	if err != nil {
		logs.Error("Init DefaultK8sBench - 1.5 fail, err: %s", err)
	}

	_, err = o.Raw(dbscript.DefaultHostImageVulnScanSql).Exec()
	if err != nil {
		logs.Error("Init DefaultHostImageVulnScanSql fail, err: %s", err)
	}

	_, err = o.Raw(dbscript.DefaultDockerVirusScanSql).Exec()
	if err != nil {
		logs.Error("Init DefaultDockerVirusScan fail, err: %s", err)
	}

	_, err = o.Raw(dbscript.DefaultHostVirusScanSql).Exec()
	if err != nil {
		logs.Error("Init DefaultHostVirusScanSql fail, err: %s", err)
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

	_, err = o.Raw(dbscript.DefaultHostImageVlunScanJobSql).Exec()
	if err != nil {
		logs.Error("Init DefaultHostImageVlunScanJobSql Job fail, err: %s", err)
	}

	_, err = o.Raw(dbscript.WarningInfoConfig).Exec()
	if err != nil {
		logs.Error("Init WarningInfo Config  fail, err: %s", err)
	}
	_, err = o.Raw(dbscript.WarningWhiteListConfig).Exec()
	if err != nil {
		logs.Error("Init WarningInfoWhiteListConfig Config  fail, err: %s", err)
	}
	_, err = o.Raw(dbscript.DefaultMailServerConfig).Exec()
	if err != nil {
		logs.Error("Init DefaultMailServer Config fail, err: %s", err)
	}

	_, err = o.Raw(dbscript.DefaultMailLogConfig).Exec()
	if err != nil {
		logs.Error("Init DefaultMailLog Config fail, err: %s", err)
	}

	_, err = o.Raw(dbscript.DefaultLDAPClientConfig).Exec()
	if err != nil {
		logs.Error("Init DefaultLDAPClient Config fail, err: %s", err)
	}

	logs.Info("Init default data end................")
}
