package models

type Registry struct {
	Id          int            `orm:"pk;auto" description:"(仓库id)"`
	Name        string         `orm:"size(64)" description:"(仓库名)"`
	Description string         `orm:"size(256)" description:"(描述/备注)"`
	Type        string         `orm:"size(32)" description:"(仓库类型)"`
	Url         string         `orm:"size(512)" dqescription:"(地址)"`
	User        string         `orm:"size(32)" description:"(用户名)"`
	Pwd         string         `orm:"size(128)" description:"(密码)"`
	ImageConfig []*ImageConfig `orm:"reverse(many);default(null)" description:"(镜像)"`
}
