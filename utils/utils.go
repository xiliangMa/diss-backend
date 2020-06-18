package utils

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"os"
	"strconv"
)

var (

	//启动模式
	Run_Mode_Dev  = "dev"
	Run_Mode_Prod = "prod"
	// 数据库 ---- 数据源(和app.config 文件中的必须对应)
	DS_Default      = "default"
	DS_Security_Log = "security_log"
	DS_Diss_Api     = "diss_api"
	// 数据库 ----数据库驱动
	DS_Driver_Mysql    = "mysql"
	DS_Driver_Postgres = "postgres"
	// 数据库 ----数据库初始化变量（default postgres）
	DS_Default_POSTGRES_ROOT_PASSWORD = "POSTGRES_ROOT_PASSWORD"
	DS_Default_POSTGRES_USER          = "POSTGRES_USER"
	DS_Default_POSTGRES_PASSWORD      = "POSTGRES_PASSWORD"
	DS_Default_POSTGRES_DB            = "POSTGRES_DB"
	DS_Default_POSTGRES_HOST          = "DEFAULT_HOST"
	DS_Default_POSTGRES_PORT          = "DEFAULT_PORT"
	// 数据库 ----数据库初始化变量（security_log postgres）
	DS_Security_Log_DB_Name = "DS_Security_Log_DB_Name"
	DS_Security_Log_User    = "DS_Security_Log_User"
	DS_Security_Log_Pwd     = "DS_Security_Log_Pwd"
	DS_Security_Log_Host    = "DS_Security_Log_Host"
	DS_Security_Log_Port    = "DS_Security_Log_Port"

	// 数据库 ----数据库初始化变量（diss_api postgres）
	DS_Diss_Api_DB_Name = "DS_Diss_Api_DB_Name"
	DS_Diss_Api_User    = "DS_Diss_Api_User"
	DS_Diss_Api_Pwd     = "DS_Diss_Api_Pwd"
	DS_Diss_Api_Host    = "DS_Diss_Api_Host"
	DS_Diss_Api_Port    = "DS_Diss_Api_Port"

	// Nats
	Nats_Prefix     = "Nats_Prefix"
	Nats_Cluster_Id = "Nats_Cluster_Id"
	Nats_Client_Id  = "Nats_Client_Id"
	Nats_Enable     = "Nats_Enable"
	Nats_Ip         = "Nats_Ip"
	Nats_Port       = "Nats_Port"
	Nats_User       = "Nats_User"
	Nats_Pwd        = "Nats_Pwd"

	// SysLog
	Syslog_Server_Url = "Syslog_Server_Url"
)

func UnitConvert(size int64) string {
	//如果字节数少于1024，则直接以B为单位，否则先除于1024，后3位因太少无意义
	if size < 1024 {
		return strconv.FormatInt(size, 10) + "B"
	} else {
		size = size / 1024
	}
	//如果原字节数除于1024之后，少于1024，则可以直接以KB作为单位
	//因为还没有到达要使用另一个单位的时候
	//接下去以此类推
	if size < 1024 {
		return strconv.FormatInt(size, 10) + "KB"
	} else {
		size = size / 1024
	}
	if size < 1024 {
		//因为如果以MB为单位的话，要保留最后1位小数，
		//因此，把此数乘以100之后再取余
		size = size * 100
		strconv.FormatInt(size, 10)
		return strconv.FormatInt((size/100), 10) + "." + strconv.FormatInt((size%100), 10) + "MB"
	} else {
		//否则如果要以GB为单位的，先除于1024再作同样的处理
		size = size * 100 / 1024
		return strconv.FormatInt((size/100), 10) + "." + strconv.FormatInt((size%100), 10) + "GB"
	}
}

func IgnoreLastInsertIdErrForPostgres(err error) error {
	msg := "LastInsertId is not supported by this driver"
	if err.Error() == msg {
		return nil
	}
	return err
}

func IgnoreQuerySeterNoRowFoundErrForPostgres(err error) error {
	msg := "<QuerySeter> no row found"
	if err.Error() == msg {
		return nil
	}
	return err
}

func GetMarkSummarySql(BMLT string) string {
	sql := "select " +
		"sum(fail_count) as fail_count, " +
		"sum(warn_count) as warn_count, " +
		"sum(info_count) as info_count, " +
		"sum(pass_count) as pass_count " +
		"from bench_mark_log " +
		"where type='" + BMLT + "'"
	return sql
}

func GetHostMarkSummarySql() string {
	sql := "select " +
		"sum(fail_count) as fail_count, " +
		"sum(warn_count) as warn_count, " +
		"sum(info_count) as info_count, " +
		"sum(pass_count) as pass_count " +
		"from bench_mark_log"
	return sql
}

/**
 * @serverUrl nats://diss:diss@111.229.167.6:4222
 */
func GetNatsServerUrl() string {
	runMode := beego.AppConfig.String("RunMode")
	envRunMode := os.Getenv("RunMode")
	if envRunMode != "" {
		runMode = envRunMode
	}
	prefix := beego.AppConfig.String("nats::Prefix")
	ip := beego.AppConfig.String("nats::Ip")
	port := beego.AppConfig.String("nats::Port")
	user := beego.AppConfig.String("nats::User")
	pwd := beego.AppConfig.String("nats::Pwd")
	serverUrl := prefix + user + `:` + pwd + `@` + ip + `:` + port
	if runMode == Run_Mode_Prod {
		prefix = os.Getenv(Nats_Prefix)
		ip = os.Getenv(Nats_Ip)
		port = os.Getenv(Nats_Port)
		user = os.Getenv(Nats_User)
		pwd = os.Getenv(Nats_Pwd)
		serverUrl = prefix + user + `:` + pwd + `@` + ip + `:` + port
	}
	logs.Info("Nats conect url: %s", serverUrl)
	return serverUrl
}

func GetNatsClientId() string {
	runMode := beego.AppConfig.String("RunMode")
	envRunMode := os.Getenv("RunMode")
	if envRunMode != "" {
		runMode = envRunMode
	}
	clientId := beego.AppConfig.String("nats::ClientId")
	if runMode == Run_Mode_Prod {
		clientId = os.Getenv(Nats_Client_Id)
	}
	return clientId
}

func GetNatsClusterId() string {
	runMode := beego.AppConfig.String("RunMode")
	envRunMode := os.Getenv("RunMode")
	if envRunMode != "" {
		runMode = envRunMode
	}
	clusterId := beego.AppConfig.String("nats::ClusterId")
	if runMode == Run_Mode_Prod {
		clusterId = os.Getenv(Nats_Cluster_Id)
	}
	return clusterId
}

func IsEnableNats() bool {
	runMode := beego.AppConfig.String("RunMode")
	envRunMode := os.Getenv("RunMode")
	if envRunMode != "" {
		runMode = envRunMode
	}
	enable := false
	if ok, _ := beego.AppConfig.Bool("nats::Enable"); ok {
		enable = true
	}
	if runMode == Run_Mode_Prod {
		ok := os.Getenv(Nats_Enable)
		if ok == "true" {
			enable = true
		}
	}
	return enable
}
