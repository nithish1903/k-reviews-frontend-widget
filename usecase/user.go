package usecase

import (
	"k-reviews-frontend-api/constant"
	"k-reviews-frontend-api/entity"
	"k-reviews-frontend-api/repository/mysql"

	"go.uber.org/zap"
)

func GetAccountById(id string) (*entity.Account, error) {
	mysqlConn := mysql.GetMySqlConn().Connection

	var account entity.Account
	result := mysqlConn.Table(constant.ACCOUNT_TABLE).Where("id = ?", id).First(&account)

	if result.Error != nil {
		zap.L().Error("Error getting account information", zap.Error(result.Error), zap.String("id:", id))
		return nil, result.Error
	}
	return &account, nil
}

func GetAccountByKey(public_key string) (*entity.Account, error) {
	mysqlConn := mysql.GetMySqlConn().Connection

	var account entity.Account
	result := mysqlConn.Table(constant.ACCOUNT_TABLE).Where("public_key = ?", public_key).First(&account)

	if result.Error != nil {
		zap.L().Error("Error getting account information", zap.Error(result.Error), zap.String("id:", public_key))
		return nil, result.Error
	}
	return &account, nil
}
