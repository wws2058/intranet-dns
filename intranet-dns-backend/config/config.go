package config

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

var GlobalConfig = new(Config)

type Config struct {
	App struct {
		Name string `ini:"name"`
		Env  string `ini:"env"`
		Port int    `ini:"port"`
	} `ini:"app"`
	Mysql struct {
		User     string `ini:"user"`
		Host     string `ini:"host"`
		Port     int    `ini:"port"`
		Passwd   string `ini:"passwd"`
		Database string `ini:"database"`
	} `ini:"mysql"`
	Redis struct {
		Host     string `ini:"host"`
		Port     int    `ini:"port"`
		Passwd   string `ini:"passwd"`
		Database int    `ini:"database"`
	} `ini:"redis"`
}

// initialize global config
func Init() {
	config, err := ini.Load("./config/config.ini")
	// config, err := ini.Load("../../config/config.ini")
	if err != nil {
		logrus.Errorf("load config.ini: %v", err)
		return
	}

	err = config.MapTo(GlobalConfig)
	if err != nil {
		logrus.Errorf("load config.ini to struct: %v", err)
		return
	}
}
