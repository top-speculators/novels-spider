package main

import (
	"gin-blog/models"
	"gin-blog/services"
	"gin-blog/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {

	// Debug 设置
	// gin.SetMode(gin.ReleaseMode) 需要时开启

	// 读取全局配置
	// 这里的配置生命周期是 main goroutine，且全局有效
	// 当子 goroutine 修改了配置，将作用于之后运行的所有其他 goroutine ，除非 main goroutine 重启
	// 当修改了配置文件，需重启 main.go 来使其生效
	utils.LoadConfig("./config.yaml")

	// sqlite 连接
	db, _ := models.Conn()
	defer db.Close()

	// 设置日志模式
	// TODO: gin 的原始 log 不区分错误日志，替换成功能更全面的 seelog
	logType := utils.GetConfig("log_type").(bool)
	logFile := utils.GetConfig("log_file").(string)
	utils.SetLogConfig(logType, logFile)

	// 注册路由并开启 http 监听
	addr := utils.GetConfig("addr").(string)
	err := services.RegisterRouter(gin.Default()).Run(addr)
	if err != nil {
		panic(err)
	}
}
