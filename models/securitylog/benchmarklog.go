package securitylog

import "github.com/astaxie/beego/orm"

type BenchMarkLog struct {
	Id            string `orm:"pk;description(基线id)"`
	BenchMarkName string `orm:"description(基线模版名称)"`
	Level         string `orm:"description(级别)"`
	ProjectName   string `orm:"description(测试项目)"`
	Result        string `orm:"null;description(测试结果)"`
	Info          string `orm:"null;description(详细信息)"`
	FailCount    string   `orm:"description(检查失败个数)"`
	WarnCount    string   `orm:"description(检查警告个数)"`
	PassCount    string   `orm:"description(检查通过个数)"`
	InfoCount    string   `orm:"description(检查提示个数)"`
}

func init() {
	orm.RegisterModel(new(BenchMarkLog))
}
