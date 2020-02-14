package models

type TenantImage struct {
	TenantId      string `orm:"description(租户id)"`
	ImageConfigId string `orm:"description(镜像id)"`
}
