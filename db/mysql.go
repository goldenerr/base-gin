package db

import (
	"github.com/basegin/base/jinzhu/gorm"
	"github.com/basegin/base/log"
	"github.com/basegin/config"
	"time"
)

var MysqlClient *gorm.DB

func InitMysql() {
	var err error
	MysqlClient, err = InitMysqlClient(MysqlConf{
		User:            config.Conf.Mysql.User,
		Password:        config.Conf.Mysql.PassWord,
		Addr:            config.Conf.Mysql.Addr,
		DataBase:        config.Conf.Mysql.DataBase,
		MaxIdleConns:    config.Conf.Mysql.MaxIdleConns,
		MaxOpenConns:    config.Conf.Mysql.MaxOpenConns,
		ConnMaxLifeTime: 3600 * time.Second,
		LogMode:         true,
	})
	if err != nil {
		log.PanicfLogger(nil, "mysql connect error: %v", err)
	}
}
