package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"
	"novels-spider/conf"
	"novels-spider/pkg/bootstrap"
	"novels-spider/services"
	"novels-spider/services/crons"
	"novels-spider/services/customers"
)

func main() {

	bootstrap.LoadLogger()

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
	addr := conf.Port
	err := services.RegisterRouter(gin.New()).Run(addr)
	if err != nil {
		logrus.Error("http 服务启动错误", err)
		return
	}
}
