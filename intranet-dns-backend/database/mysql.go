package database

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wws2058/intranet-dns/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// db connection
func initSql() {
	conf := config.GlobalConfig.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.User, conf.Passwd, conf.Host, conf.Port, conf.Database)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.WithField("mysql", "init").Error(err)
		os.Exit(1)
	}
	if config.GlobalConfig.App.Env == "dev" {
		DB = DB.Debug()
	}
	db, _ := DB.DB()
	db.SetConnMaxIdleTime(time.Hour)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	logrus.WithField("mysql", "init").Info("succeed")
}
