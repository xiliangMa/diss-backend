package dissapi

import (
	"github.com/xiliangMa/diss-backend/models"
)

func GetAccountByUser(accountUsers *models.AccountUsers) (error, string) {
	err, accoount := accountUsers.GetAccountByUser()
	if err != nil {
		return nil, accoount.AccountName
	}
	return err, ""
}