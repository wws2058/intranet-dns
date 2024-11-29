package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/tswcbyy1107/dns-service/apis"
	"github.com/tswcbyy1107/dns-service/config"
	"github.com/tswcbyy1107/dns-service/middleware"
	"github.com/tswcbyy1107/dns-service/models"

	_ "github.com/tswcbyy1107/dns-service/docs"
)

func InitRouter() *http.Server {
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	r.Use(middleware.LogHandler())
	r.Use(middleware.Auth())
	r.Use(middleware.ApiLimiter())

	// gin swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// load apis
	apis.LoadPingApis(r)
	apis.LoadSysApis(r)
	apis.LoadCronjobApis(r)

	// clean unused apis in db
	ginApis := []models.Api{}
	for _, router := range r.Routes() {
		ginApis = append(ginApis, models.Api{
			Path:   router.Path,
			Method: router.Method,
		})
	}
	models.ApisCleanInDB(ginApis)

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%v", config.GlobalConfig.App.Port),
		Handler:        r,
		ReadTimeout:    time.Second * 30,
		WriteTimeout:   time.Second * 30,
		MaxHeaderBytes: 1 << 20,
	}
	return srv
}
