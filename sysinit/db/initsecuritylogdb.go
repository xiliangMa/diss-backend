package db

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
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
	//force, _ := web.AppConfig.Bool("Force")

	//连接名称
	//dbAlias := web.AppConfig.String(DS + "::Alias")
	//数据库名称
	dbName, _ := web.AppConfig.String(DSAlias + "::DBName")
	//数据库连接用户名
	dbUser, _ := web.AppConfig.String(DSAlias + "::User")
	//数据库连接密码
	dbPwd, _ := web.AppConfig.String(DSAlias + "::Pwd")
	//数据库IP（域名）
	dbHost, _ := web.AppConfig.String(DSAlias + "::Host")
	//端口
	port, _ := web.AppConfig.String(DSAlias + "::Port")
	// postgres
	DS := fmt.Sprintf("%s%s%s%s%s%s", "host="+dbHost, " port="+port, " user="+dbUser, " password="+dbPwd, " dbname="+dbName, " sslmode=disable")
	orm.RegisterDriver(driver, orm.DRPostgres)
	err := orm.RegisterDataBase(DSAlias, driver, DS)
	if err != nil {
		logs.Error("DB Register fail, >>> DSAlias: %s <<<, Err: %s", DSAlias, err)
	}

	//创建数据表
	this.CreateOrUpdateDb()
}

func (this *SecurityLogDb) CreateOrUpdateDb() {
	dbName := utils.DS_Security_Log
	o := orm.NewOrmUsingDB(dbName)
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
