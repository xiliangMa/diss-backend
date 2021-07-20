package models

// 客户端模块控制模型(虚拟模型）
type ClientModuleControl struct {
	ModuleName string
	Enable     bool
	Status     string
}

// 通过客户端更新资产数据(虚拟模型）
type UpdateAssets struct {
	GlobalRefresh bool
	Tags          *map[string]string
}
