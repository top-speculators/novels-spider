package main

import (
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"novels-spider/pkg/bootstrap"
	"novels-spider/pkg/helpers"
	"novels-spider/services"
	"novels-spider/services/crons"
	"novels-spider/services/customers"
)

func main() {

	// bootstrap 的顺序不可打乱

	bootstrap.LoadLogger()
	defer bootstrap.LoggerFlush()

	bootstrap.LoadLocalConfig()

	bootstrap.SetDebugMode()

	bootstrap.LoadDBConnections()
	defer bootstrap.DBClose()

	bootstrap.LoadMQConnection()
	defer bootstrap.MQClose()

	// cron
	go crons.StartListen()

	// customers
	customers.StartListen()

	// web 服务
	addr := helpers.GetConfig("port").(string)
	err := services.RegisterRouter(gin.New()).Run(addr)
	if err != nil {
		_ = seelog.Critical("http 服务启动错误", err)
		return
	}
}
