package system

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"regexp"
	"strings"
)

type WariningFilterService struct {
	Info string
	Rule string
}

func (this *WariningFilterService) CheckFromWhiteListItem() bool {
	rule := strings.Replace(this.Rule, models.WarnWhiteListCnTrans_Node[0], models.WarnWhiteListCnTrans_Node[1], 1)
	rule = strings.Replace(rule, models.WarnWhiteListCnTrans_ContainerId[0], models.WarnWhiteListCnTrans_ContainerId[1], 1)
	rule = strings.Replace(rule, models.WarnWhiteListCnTrans_ContainerName[0], models.WarnWhiteListCnTrans_ContainerName[1], 1)
	rule = strings.Replace(rule, models.WarnWhiteListCnTrans_CmdLine[0], models.WarnWhiteListCnTrans_CmdLine[1], 1)

	warnRule := []byte(rule)
	rulelines := map[string]string{}
	json.Unmarshal(warnRule, &rulelines)

	for rulekey, ruleitem := range rulelines {
		if ruleitem != "" {
			val := strings.ReplaceAll(ruleitem, ".", "\\.")
			whitelistRegex := regexp.MustCompile(`"` + rulekey + `":".*` + val + `.*"`)
			match := whitelistRegex.Match([]byte(this.Info))
			if !match {
				return false
			}
		}
	}

	return true
}

func (this *WariningFilterService) FilterWarnWhiteList(warningInfo models.WarningInfo) bool {
	whiteList := models.WarningWhiteList{}
	whiteList.Enabled = true
	avaiWhiteList, _, err := whiteList.WhiteList(0, 0)
	if err != nil {
		logs.Error("Get WarningWhiteList to Apply failed, code: %d, err: %s", utils.GetWarningWhiteListErr, err.Error())
	}

	if len(avaiWhiteList) > 0 {
		for _, whiteItem := range avaiWhiteList {
			if whiteItem.WarningInfoType == warningInfo.Type && whiteItem.WarningInfoName == warningInfo.Name {
				warnFilterservice := new(WariningFilterService)
				warnFilterservice.Info = warningInfo.Info
				warnFilterservice.Rule = whiteItem.Rule
				checkstatus := warnFilterservice.CheckFromWhiteListItem()
				if checkstatus {
					return true
				}
			}
		}
	}
	return false
}
