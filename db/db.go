package db

import (
	"database/sql"
	"fmt"
	"k8s-platform/config"
	"time"

	"github.com/wonderivan/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	isInit bool
	DB     *gorm.DB
	sqlDB  *sql.DB
	err    error
)

func Init() {
	// 已经初始化直接返回
	if isInit {
		return
	}

	// 初始化连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DbUser, config.DbPwd, config.DbHost, config.DbPort, config.DbName)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		}})

	if err != nil {
		panic("db连接失败" + err.Error())
	}

	// 连接池
	sqlDB, err = DB.DB()
	if err != nil {
		panic("db连接池调整失败" + err.Error())
	}
	sqlDB.SetMaxIdleConns(config.DbMaxIdles)
	sqlDB.SetMaxOpenConns(config.DbMaxConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.DbMaxLifetime))

	// 修改flag
	isInit = true
	logger.Info("db连接成功")

}

// 关闭连接池
func Close() error {
	return sqlDB.Close()
}
