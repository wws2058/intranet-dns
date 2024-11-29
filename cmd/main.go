package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tswcbyy1107/dns-service/config"
	"github.com/tswcbyy1107/dns-service/database"
	"github.com/tswcbyy1107/dns-service/models"
	"github.com/tswcbyy1107/dns-service/router"
	"github.com/tswcbyy1107/dns-service/service/cronjob"
	"github.com/tswcbyy1107/dns-service/utils"
)

// @title        sre dns service backend demo
// @version      1.0
// @description  simple intranet dns management system, dns crud operation with centos bind9
// @host         localhost:16789
// @accept       json
// @produce      json
// @schemes      http
func main() {
	logrus.WithField("server", "staring").Info("gin")

	srv := router.InitRouter()
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithField("server", "failed").Error(err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	// ctrl c + kill
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	// graceful terminated
	timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := srv.Shutdown(timeout); err != nil {
		logrus.WithField("server", "graceful terminating err").Error(err)
	}
	logrus.WithField("server", "graceful terminated").Info("")
}

// initialize
func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: utils.DefaultTimeFormat,
	})
	logrus.SetReportCaller(true)

	config.Init()
	database.InitDB()
	models.AutoMigrate()
	cronjob.InitCronJob()
}
