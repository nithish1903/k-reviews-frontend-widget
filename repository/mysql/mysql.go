package mysql

import (
	"fmt"
	"os"

	"github.com/matryer/resync"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var onceMysql resync.Once
var mysqlConn *gorm.DB

type MysqlCon struct {
	Connection *gorm.DB
}

func GetMySqlConn() *MysqlCon {
	onceMysql.Do(func() {
		zap.L().Info("Creating MySQl connection")
		userName := os.Getenv("MYSQL_USERNAME")
		password := os.Getenv("MYSQL_PASSWORD")
		dbHost := os.Getenv("MYSQL_HOST")
		dbName := os.Getenv("MYSQL_DB_NAME")
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", userName, password, dbHost, dbName)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			zap.L().Error("Could not create MySql connection", zap.Any("err", err))
		}
		zap.L().Info("Created MySQl DB connection")
		mysqlConn = db
	})
	return &MysqlCon{Connection: mysqlConn}
}
