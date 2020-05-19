package models

import "time"

type LicenseFile struct {
	Id                string           `orm:"pk;" description:"(license file id)"`
	ProductName       string           `orm:"" description:"(产品名称)"`
	CustomerName      string           `orm:"" description:"(许可对象)"`
	LicenseType       int              `orm:"" description:"(授权类型 0测试 1正式)"`
	LicenseUuid       string           `orm:"" description:"(序列号)"`
	LicenseBuyAt      time.Time        `orm:"" description:"(授权购买时间)"`
	LicenseActiveAt   time.Time        `orm:"" description:"(激活时间)"`
	LicenseModule     []*LicenseModule `orm:"reverse(many);null" description:"(授权的模块)"`
	LicenseModuleJson string           `orm:"" description:"(授权的模块Json 格式版)"`
}

type LicenseModule struct {
	Id              string       `orm:"pk;" description:"(license module id)"`
	LincenseFile    *LicenseFile `orm:"rel(fk);null;" description:"(license file)"`
	ModuleCode      string       `orm:"" description:"(授权模块)"`
	LicenseCount    string       `orm:"" description:"(授权模块数量)"`
	LicenseExpireAt time.Time    `orm:"" description:"(授权结束时间)"`
}
