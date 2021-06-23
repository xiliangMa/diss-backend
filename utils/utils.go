package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/shirou/gopsutil/host"
	"strconv"
	"strings"
)

var (
	// 数据库 ---- 数据源(和app.config 文件中的必须对应)
	DS_Default      = "default"
	DS_Security_Log = "security_log"
	// 数据库 ----数据库驱动
	DS_Driver_Postgres = "postgres"
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

func GetHostMarkSummarySql(hostId string) string {
	sql := "select " +
		"sum(fail_count) as fail_count, " +
		"sum(warn_count) as warn_count, " +
		"sum(info_count) as info_count, " +
		"sum(pass_count) as pass_count " +
		"from bench_mark_log " +
		"where host_id ='" + hostId + "'"
	return sql
}

func GetDBCountSql(hostId string) string {
	sql := `SELECT 
			mysql.mysql_count, 
			oracle.oracle_count, 
			redis.redis_count, 
			postgres.postgres_count, 
			mongodb.mongodb_count, 
			memcache.memcache_count, 
			db2.db2_count, 
			hbase.hbase_count 
			from 
			(SELECT count(id) as mysql_count from ` + ImageConfig + ` WHERE name ilike '%mysql%') as mysql ,
			(SELECT count(id) as oracle_count from ` + ImageConfig + ` WHERE name ilike '%oracle%') as oracle,
			(SELECT count(id) as redis_count from ` + ImageConfig + ` WHERE name ilike '%redis%') as redis,
			(SELECT count(id) as postgres_count from ` + ImageConfig + ` WHERE name ilike '%postgres%') as postgres,
			(SELECT count(id) as mongodb_count from ` + ImageConfig + ` WHERE name ilike '%mongodb%') as mongodb,
			(SELECT count(id) as memcache_count from ` + ImageConfig + ` WHERE name ilike '%memcache%') as memcache,
			(SELECT count(id) as db2_count from ` + ImageConfig + ` WHERE name ilike '%db2%') as db2,
			(SELECT count(id) as hbase_count from ` + ImageConfig + ` WHERE name ilike '%hbase%') as hbase;`
	if hostId == "" {
		return sql
	}
	sqlByHostId := `SELECT 
			mysql.mysql_count, 
			oracle.oracle_count, 
			redis.redis_count, 
			postgres.postgres_count, 
			mongodb.mongodb_count, 
			memcache.memcache_count, 
			db2.db2_count, 
			hbase.hbase_count 
			from 
			(SELECT count(id) as mysql_count from ` + ImageConfig + ` WHERE name ilike '%mysql%' and host_id = ?) as mysql ,
			(SELECT count(id) as oracle_count from ` + ImageConfig + ` WHERE name ilike '%oracle%' and host_id = ?) as oracle,
			(SELECT count(id) as redis_count from ` + ImageConfig + ` WHERE name ilike '%redis%' and host_id = ?) as redis,
			(SELECT count(id) as postgres_count from ` + ImageConfig + ` WHERE name ilike '%postgres%' and host_id = ?) as postgres,
			(SELECT count(id) as mongodb_count from ` + ImageConfig + ` WHERE name ilike '%mongodb%' and host_id = ?) as mongodb,
			(SELECT count(id) as memcache_count from ` + ImageConfig + ` WHERE name ilike '%memcache%' and host_id = ?) as memcache,
			(SELECT count(id) as db2_count from ` + ImageConfig + ` WHERE name ilike '%db2%' and host_id = ?) as db2,
			(SELECT count(id) as hbase_count from ` + ImageConfig + ` WHERE name ilike '%hbase%' and host_id = ?) as hbase;`

	return sqlByHostId
}

/**
 * @serverUrl nats://diss:diss@111.229.167.6:4222
 */
func GetNatsServerUrl() string {
	prefix := beego.AppConfig.String("nats::Prefix")
	ip := beego.AppConfig.String("nats::Ip")
	port := beego.AppConfig.String("nats::Port")
	user := beego.AppConfig.String("nats::User")
	pwd := beego.AppConfig.String("nats::Pwd")
	serverUrl := prefix + user + `:` + pwd + `@` + ip + `:` + port
	logs.Info("Nats conect url: %s", serverUrl)
	return serverUrl
}

func GetHTTPSAddr() string {
	HTTPSAddr := beego.AppConfig.String("HTTPSAddr")
	return HTTPSAddr
}

func GetHTTPSPort() string {
	HTTPSPort := beego.AppConfig.String("HTTPSPort")
	return HTTPSPort
}

func GetNatsClientId() string {
	clientId := beego.AppConfig.String("nats::ClientId")
	return clientId
}

func GetNatsClusterId() string {
	clusterId := beego.AppConfig.String("nats::ClusterId")
	return clusterId
}

func IsEnableNats() bool {
	enable := false
	if ok, _ := beego.AppConfig.Bool("nats::Enable"); ok {
		enable = true
	}
	return enable
}

func GetScopeYml() string {
	scopeYml := beego.AppConfig.String("k8s::ScopeYml")
	return scopeYml
}

func GetKubeScanYml() string {
	kubescanYml := beego.AppConfig.String("k8s::KubeScanYml")
	return kubescanYml
}

func GetKubeScanEnvDisPatchUrlKey() string {
	KubeScanEnvDisPatchUrlKey := beego.AppConfig.String("k8s::KubeScanEnvDisPatchUrlKey")
	return KubeScanEnvDisPatchUrlKey
}

func GetKubeScanEnvTaskIdKey() string {
	KubeScanEnvTaskIdKey := beego.AppConfig.String("k8s::KubeScanEnvTaskIdKey")
	return KubeScanEnvTaskIdKey
}

func GetKubeScanEnvClusterIdKey() string {
	KubeScanEnvClusterIdKey := beego.AppConfig.String("k8s::KubeScanEnvClusterIdKey")
	return KubeScanEnvClusterIdKey
}

func GetKubeScanReportApi() string {
	KubeScanReportApi := beego.AppConfig.String("k8s::KubeScanReportApi")
	return KubeScanReportApi
}

func GetKubeScanReportUrl() string {
	sechemes := "https://"
	ip := GetSystemIP()
	port := GetHTTPSPort()
	repotUri := fmt.Sprintf("%s%s:%s%s", sechemes, ip, port, GetKubeScanReportApi())
	return repotUri
}

func GetScopeNameSpace() string {
	scopeNameSpace := beego.AppConfig.String("k8s::ScopeNameSpace")
	return scopeNameSpace
}

func GetScopeAppName() string {
	scopeAppName := beego.AppConfig.String("k8s::ScopeAppName")
	return scopeAppName
}

func ConvertType(from interface{}, to interface{}) error {
	data, err := json.Marshal(from)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, to)
	if err != nil {
		return err
	}
	return nil
}

func GetLogoPath() string {
	return beego.AppConfig.String("system::LogoPath")
}

func GetLogoName() string {
	return beego.AppConfig.String("system::NewLogoName")
}

func GetDefaultLogoName() string {
	return beego.AppConfig.String("system::DefaultLogoName")
}

func GetLogoUrl() string {
	return beego.AppConfig.String("system::LogoUrl")
}

func GetWarnWhitelistPath() string {
	return beego.AppConfig.String("system::WarnWhitelistPath")
}

func GetVulnDbPath() string {
	return beego.AppConfig.String("system::VulnDbPath")
}

func GetVulnDbUrl() string {
	return beego.AppConfig.String("system::VulnDbUrl")
}

func GetProbeDriverPath() string {
	return beego.AppConfig.String("system::ProbeDriverPath")
}

func GetProbeDriverUrl() string {
	return beego.AppConfig.String("system::ProbeDriverUrl")
}

func GetVirusPath() string {
	return beego.AppConfig.String("system::VirusPath")
}

func GetVirusrUrl() string {
	return beego.AppConfig.String("system::VirusUrl")
}

func GetPublicPath() string {
	return beego.AppConfig.String("system::PublicPath")
}

func GetPublicUrl() string {
	return beego.AppConfig.String("system::PublicUrl")
}

func GetSystemIP() string {
	IP := beego.AppConfig.String("system::IP")
	return IP
}

func FormatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.1fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.1fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.1fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.1fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.1fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else {
		//if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.1fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

func InitTable(table string, sql string) (pre string) {
	pre = `DO
		$$
		DECLARE
			tableName varchar :='` + table + `';
		BEGIN
		IF NOT EXISTS (SELECT * FROM information_schema.columns WHERE table_name=tableName) THEN
		` + sql + `
		RAISE NOTICE 'Create Table % %', tableName, 'Success';
		ELSE
		RAISE NOTICE 'Already exists';
		END IF;
		END
		$$`
	return pre
}

func ToIndentJSON(obj interface{}) (string, error) {
	bs, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = json.Indent(buf, bs, "", "\t")
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func GetHostInfo() *host.InfoStat {
	info, _ := host.Info()
	return info
}

// 这个方法在Name字段加个前缀role_
func GetRoleString(s string) string {
	if strings.HasPrefix(s, "role_") {
		return s
	}
	return fmt.Sprintf("role_%s", s)
}

func MD5(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}
