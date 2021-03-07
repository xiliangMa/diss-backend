package models

type DisposalMode struct {
	WarningInfo      WarningInfo
	WarningWhiteList WarningWhiteList
	Action           string `orm:"-" description:"(处理方式：isolation、pause、stop、kill)"`
}

type DisposalModeInterface interface {
}
