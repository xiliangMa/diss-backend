package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type NameSpace struct {
	Id             string `orm:"pk;" description:"(命名空间id)"`
	Name           string `orm:"size(32)" description:"(命名空间)"`
	ClusterId      string `orm:"size(64);" description:"(集群id)"`
	ClusterName    string `orm:"size(32)" description:"(集群名)"`
	AccountName    string `orm:"size(32);default(admin)" description:"(租户)"`
	SyncCheckPoint int64  `orm:"default(0);" description:"(同步检查点)"`
	Force          bool   `orm:"-" description:"(强制)"`
	KMetaData      string `orm:"" description:"(源数据)"`
	KSpec          string `orm:"" description:"(Spec数据)"`
	KStatus        string `orm:"" description:"(状态数据)"`
}

type NameSpaceInterface interface {
	Add() Result
	Delete() Result
	List() Result
	BindAccount() Result
	UnBindAccount() Result
	ListByAccountGroupByClusterId() (error, []string)
	EmptyDirtyData() error
}

func (this *NameSpace) Add(update bool) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	var dbNS []*NameSpace
	var err error
	cond := orm.NewCondition()
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	count, err := o.QueryTable(utils.NameSpace).SetCond(cond).All(&dbNS)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetNameSpaceErr
		logs.Error("Get NameSpace failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	if count != 0 {
		if update {
			// 同步更新k8s数据
			this.AccountName = dbNS[0].AccountName
			this.Update()
		} else {
			ResultData.Message = "NameSpaceExistErr"
			ResultData.Code = utils.NameSpaceExistErr
			logs.Error("Add NameSpace failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
			return ResultData
		}
	} else {
		_, err = o.Insert(this)
		if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
			ResultData.Message = err.Error()
			ResultData.Code = utils.AddNameSpaceErr
			logs.Error("Add NameSpace failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
			return ResultData
		}
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *NameSpace) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var nameSpaceList []*NameSpace
	var ResultData Result
	var err error
	cond := orm.NewCondition()

	if this.Name != "" {
		cond = cond.And("name__contains", this.Name)
	}
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}

	if this.ClusterId != "" {
		cond = cond.And("cluster_id", this.ClusterId)
	}

	if this.ClusterName != "" {
		cond = cond.And("cluster_name__contains", this.ClusterName)
	}

	if this.AccountName != "" && this.AccountName != Account_Admin {
		cond = cond.And("account_name", this.AccountName)
	}

	_, err = o.QueryTable(utils.NameSpace).SetCond(cond).Limit(limit, from).All(&nameSpaceList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetClusterErr
		logs.Error("Get NameSpace List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.NameSpace).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = nameSpaceList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *NameSpace) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditNameSpaceErr
		logs.Error("Update NameSpace: %s failed, code: %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *NameSpace) UnBindAccount() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	err := o.Begin()
	params := orm.Params{"account_name": Account_Admin}
	cond := orm.NewCondition()
	cond = cond.And("id", this.Id)

	_, err = o.QueryTable(utils.NameSpace).SetCond(cond).Update(params)
	if err != nil {
		o.Rollback()
		ResultData.Message = err.Error()
		ResultData.Code = utils.UnBindNameSpaceErr
		logs.Error("UnBind NameSpace: %s failed, code: %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
		return ResultData
	}
	// 解除 pod 和 pod 下容器的 租户绑定关系
	_, err = o.QueryTable(utils.Pod).Filter("name_space_name", this.Name).Update(params)
	if err != nil {
		o.Rollback()
		ResultData.Message = err.Error()
		ResultData.Code = utils.UnBindPodErr
		logs.Error("UnBind Pod failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	_, err = o.QueryTable(utils.ContainerConfig).Filter("name_space_name", this.Name).Update(params)
	if err != nil {
		o.Rollback()
		ResultData.Message = err.Error()
		ResultData.Code = utils.UnBindContainerErr
		logs.Error("UnBind Container failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	o.Commit()
	ResultData.Code = http.StatusOK
	ResultData.Data = nil
	return ResultData
}

func (this *NameSpace) BindAccount() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()
	cond = cond.And("id", this.Id)
	dbList := []*NameSpace{}

	o.QueryTable(utils.NameSpace).SetCond(cond).All(&dbList)
	if dbList == nil {
		ResultData.Message = "NoNameSpacedErr"
		ResultData.Code = utils.NoNameSpacedErr
		logs.Error("NameSpace: %s not found, code: %d, err: %s", dbList[0].Name, ResultData.Code, ResultData.Message)
		return ResultData
	}

	if this.Force {
		dbList[0].AccountName = this.AccountName
	} else {
		if dbList[0].AccountName == "" || dbList[0].AccountName == Account_Admin {
			dbList[0].AccountName = this.AccountName
		} else {
			if this.Force != true {
				o.Rollback()
				ResultData.Message = "IsBindErr"
				ResultData.Code = utils.IsBindErr
				logs.Error("NameSpace: %s bind accout failed, code: %d, err: %s", dbList[0].Name, ResultData.Code, ResultData.Message)
				return ResultData
			}
		}
	}
	dbList[0].Force = this.Force

	_, err := o.Update(dbList[0])
	if err != nil {
		o.Rollback()
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditNameSpaceErr
		logs.Error("Bind NameSpace: %s, AccountNamecode: %s, fail: %d, err: %s", dbList[0].Name, this.AccountName, ResultData.Code, ResultData.Message)
		return ResultData
	}
	// 更新 pod 和 pod 下容器的 租户绑定关系
	params := orm.Params{"account_name": this.AccountName}
	_, err = o.QueryTable(utils.Pod).Filter("name_space_name", this.Name).Update(params)
	if err != nil {
		o.Rollback()
		ResultData.Message = err.Error()
		ResultData.Code = utils.UnBindPodErr
		logs.Error("UnBind Pod failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	_, err = o.QueryTable(utils.ContainerConfig).Filter("name_space_name", this.Name).Update(params)
	if err != nil {
		o.Rollback()
		ResultData.Message = err.Error()
		ResultData.Code = utils.UnBindContainerErr
		logs.Error("UnBind Container failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	o.Commit()
	ResultData.Code = http.StatusOK
	ResultData.Data = dbList[0]
	return ResultData
}

func (this *NameSpace) ListByAccountGroupByClusterId() (error, []string) {
	o := orm.NewOrm()

	o.Using(utils.DS_Default)
	var cIds []string
	cond := orm.NewCondition()

	if this.AccountName != "" {
		cond = cond.And("account_name", this.AccountName)
	}
	_, err := o.Raw("select cluster_id from "+utils.NameSpace+" where account_name = ? group by cluster_id", this.AccountName).QueryRows(&cIds)
	//_, err = o.QueryTable(utils.NameSpace).SetCond(cond).GroupBy("cluster_id").All(&nameSpaceList)

	if err != nil {
		logs.Error("Get NameSpace List failed, code: %d, err: %s", err.Error(), utils.GetClusterErr)
		return err, nil
	}
	return nil, cIds
}

func (this *NameSpace) EmptyDirtyData() error {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	_, err := o.Raw("delete from "+utils.NameSpace+" where cluster_id = ? and sync_check_point != ? ", this.ClusterId, this.SyncCheckPoint).Exec()
	if err != nil {
		logs.Error("Empty Dirty Data failed,  model: %s, code: %d, err: %s", utils.NameSpace, utils.EmptyDirtyDataNameSpaceErr, err.Error())
	} else {
		logs.Info("Empty Dirty Data success,  model: %s, ", utils.NameSpace)
	}
	return err
}

func (this *NameSpace) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if this.ClusterId != "" {
		cond = cond.And("cluster_id", this.ClusterId)
	}
	_, err := o.QueryTable(utils.NameSpace).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteNameSpaceErr
		logs.Error("Delete NameSpace failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}
