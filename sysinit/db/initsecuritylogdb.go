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

	table := make(map[string]string)
	table[utils.CmdHistory] = timescaledb.Tab_Create_CmdHistory
	table[utils.DockerEvent] = timescaledb.Tab_Create_DockerEvent
	table[utils.UserEvent] = timescaledb.Tab_Create_UserEvent
	table[utils.TaskLog] = timescaledb.Tab_Create_TaskLog
	table[utils.VirusScan] = timescaledb.Tab_Create_VirusScan
	table[utils.VirusRecord] = timescaledb.Tab_Create_VirusRecord
	table[utils.ImageDetail] = timescaledb.Tab_Create_ImageDetail
	table[utils.SensitiveInfo] = timescaledb.Tab_Create_SensitiveInfo
	table[utils.WarningInfo] = timescaledb.Tab_Create_WarningInfo
	table[utils.HostPackage] = timescaledb.Tab_Create_HostPackage

	for k, v := range table {
		logs.Info("Create tab Table: %s >>>>>>>>>>>>>>>>", k)
		sql := utils.InitTable(k, v)
		_, err := o.Raw(sql).Exec()
		if err != nil {
			logs.Error("Create tab Table: %s fail, err: %s ", k, err)
		}
	}
}
