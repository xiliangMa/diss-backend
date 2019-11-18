package sysinit

import (
	"fmt"
	"github.com/astaxie/beego"
	_ "github.com/xiliangMa/diss-backend/models"
	"os"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func initDB() {

	dbType := beego.AppConfig.String("db::Type")

	//连接名称
	dbAlias := beego.AppConfig.String(dbType + "::Alias")

	//数据库名称
	//dbName := beego.AppConfig.String(dbType + "::Name")
	dbName := os.Getenv("MARIADB_DATABASE")

	//数据库连接用户名
	//dbUser := beego.AppConfig.String(dbType + "::User")
	dbUser := os.Getenv("MARIADB_USER")

	//数据库连接用户名
	//dbPwd := beego.AppConfig.String(dbType + "::Pwd")
	dbPwd := os.Getenv("MARIADB_PASSWORD")

	//数据库IP（域名）
	//dbHost := beego.AppConfig.String(dbType + "::Host")
	dbHost := os.Getenv("MARIADB_HOST")

	//数据库端口
	//dbPort := beego.AppConfig.String(dbType + "::Port")

	datasource := fmt.Sprintf("%s%s%s%s%s%s%s%s", dbUser, ":", dbPwd, "@tcp(", dbHost, ":3306)/", dbName, "?charset=utf8")

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase(dbAlias, "mysql", datasource)

	// local
	//orm.RegisterDataBase("default", "mysql", "root:abc123@tcp(127.0.0.1:3306)/uranus_local?charset=utf8")

	//如果是开发模式，则显示命令信息
	isDev := (beego.AppConfig.String("RunMode") == "dev")

	// true: drop table 后再建表
	force := false

	//auto create db
	orm.RunSyncdb("default", force, isDev)
}
