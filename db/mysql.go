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
	var logMode bool
	if config.Conf.LogLevel == "debug" {
		logMode = true
	}
	MysqlClient, err = InitMysqlClient(MysqlConf{
		User:            config.Conf.Mysql.User,
		Password:        config.Conf.Mysql.PassWord,
		Addr:            config.Conf.Mysql.Addr,
		DataBase:        config.Conf.Mysql.DataBase,
		MaxIdleConns:    config.Conf.Mysql.MaxIdleConns,
		MaxOpenConns:    config.Conf.Mysql.MaxOpenConns,
		ConnMaxLifeTime: 3600 * time.Second,
		LogMode:         logMode,
	})
	if err != nil {
		log.PanicfLogger(nil, "mysql connect error: %v", err)
	}

	ConfigChange()
}

func ConfigChange() {
	go func() {
		for {
			select {
			case <-config.ChangeConfig:
				InitMysql()
			}
		}
	}()
}
