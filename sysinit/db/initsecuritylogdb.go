package db

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
	"github.com/xiliangMa/diss-backend/sysinit/dbscript/timescaledb"
	"github.com/xiliangMa/diss-backend/utils"
)

type SecurityLogDb struct {
}

func (this *SecurityLogDb) InitDB() {
	driver := utils.DS_Driver_Postgres
	DSAlias := utils.DS_Security_Log
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
	if err != nil {
		logs.Error("DB Register fail, >>> DSAlias: %s <<<, Err: %s", DSAlias, err)
		return
	}
	orm.SetMaxIdleConns(DSAlias, 2)
	orm.SetMaxOpenConns(DSAlias, 100)

	//创建数据表
	this.CreateOrUpdateDb()
}

func (this *SecurityLogDb) CreateOrUpdateDb() {
	dbName := utils.DS_Security_Log
	o := orm.NewOrm()
	o.Using(dbName)
	logs.Info("Create Or Update Db: %s data start................", dbName)

	logs.Info("Create tab Table: %s >>>>>>>>>>>>>>>>", utils.CmdHistory)
	_, err := o.Raw(timescaledb.Tab_Create_CmdHistory).Exec()
	if err != nil {
		logs.Error("Create tab Table: %s fail, err: %s ", utils.CmdHistory, err)
	}

	logs.Info("Create tab Table: %s >>>>>>>>>>>>>>>>", utils.DockerEvent)
	_, err = o.Raw(timescaledb.Tab_Create_DockerEvent).Exec()
	if err != nil {
		logs.Error("Create tab Table: %s fail, err: %s ", utils.DockerEvent, err)
	}

	logs.Info("Create tab Table: %s >>>>>>>>>>>>>>>>", utils.UserEvent)
	_, err = o.Raw(timescaledb.Tab_Create_UserEvent).Exec()
	if err != nil {
		logs.Error("Create tab Table: %s fail, err: %s ", utils.UserEvent, err)
	}

	logs.Info("Create tab Table: %s >>>>>>>>>>>>>>>>", utils.TaskLog)
	_, err = o.Raw(timescaledb.Tab_Create_TaskLog).Exec()
	if err != nil {
		logs.Error("Create tab Table: %s fail, err: %s ", utils.TaskLog, err)
	}

	logs.Info("Create tab Table: %s >>>>>>>>>>>>>>>>", utils.VirusScan)
	_, err = o.Raw(timescaledb.Tab_Create_VirusScan).Exec()
	if err != nil {
		logs.Error("Create tab Table: %s fail, err: %s ", utils.VirusScan, err)
	}

	logs.Info("Create tab Table: %s >>>>>>>>>>>>>>>>", utils.VirusRecord)
	_, err = o.Raw(timescaledb.Tab_Create_VirusRecord).Exec()
	if err != nil {
		logs.Error("Create tab Table: %s fail, err: %s ", utils.VirusRecord, err)
	}

	logs.Info("Create tab Table: %s >>>>>>>>>>>>>>>>", utils.ImageDetail)
	_, err = o.Raw(timescaledb.Tab_Create_ImageDetail).Exec()
	if err != nil {
		logs.Error("Create tab Table: %s fail, err: %s ", utils.ImageDetail, err)
	}

	logs.Info("Create tab Table: %s >>>>>>>>>>>>>>>>", utils.SensitiveInfo)
	_, err = o.Raw(timescaledb.Tab_Create_SensitiveInfo).Exec()
	if err != nil {
		logs.Error("Create tab Table: %s fail, err: %s ", utils.SensitiveInfo, err)
	}

	logs.Info("Create tab Table: %s >>>>>>>>>>>>>>>>", utils.WarningInfo)
	_, err = o.Raw(timescaledb.Tab_Create_WarningInfo).Exec()
	if err != nil {
		logs.Error("Create tab Table: %s fail, err: %s ", utils.WarningInfo, err)
	}

	logs.Info("Create tab Table: %s >>>>>>>>>>>>>>>>", utils.HostPackage)
	_, err = o.Raw(timescaledb.Tab_Create_HostPackage).Exec()
	if err != nil {
		logs.Error("Create tab Table: %s fail, err: %s ", utils.HostPackage, err)
	}
}
