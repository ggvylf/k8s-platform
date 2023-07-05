package db

import (
	"fmt"
	"k8s-platform/config"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/wonderivan/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm/schema"
)

var (
	isInit bool
	Db     *gorm.DB
	err    error
)

func Init() {
	// 已经初始化直接返回
	if isInit {
		return
	}

	// 初始化连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DbUser, config.DbPwd, config.DbHost, config.DbPort, config.DbName)
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{NameingStrategy: schema.NamingStrategy{
		SingularTable: true,
	}})

	if err != nil {
		panic("db连接失败" + err.Error())
	}

	// 调整连接池
	Db.DB().SetMaxIdleConns(config.DbMaxIdles)
	Db.DB().SetMaxOpenConns(config.DbMaxConns)
	Db.DB().SetConnMaxLifetime(time.Duration(config.DbMaxLifetime))

	// 修改flag
	isInit = true
	logger.Info("db连接成功")

}
