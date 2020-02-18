package main

import (
	"gin-blog/models"
	"gin-blog/services"
	"gin-blog/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {

	// 读取全局配置
	// 这里的配置生命周期是 main goroutine，且全局有效
	// 当子 goroutine 修改了配置，将作用于之后开启的所有其他 goroutine ，除非 main goroutine 重启
	// 当修改了配置文件，需重启项目来使其生效
	// 注意这个配置，无法生效在其他包的 init 函数，因为那时候还未载入
	err := utils.LoadConfig("./config.yaml")
	if err != nil {
		panic(err)
	}

	// debug 设置
	debug := utils.GetConfig("debug").(bool)
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// sqlite3 连接
	db, _ := models.Conn()
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()

	// 设置日志模式
	// TODO: gin 的原始 log 不区分错误日志，替换成功能更全面的 seelog
	logType := utils.GetConfig("log_type").(bool)
	logFile := utils.GetConfig("log_file").(string)
	utils.SetLogConfig(logType, logFile)

	// 注册路由并开启 http 监听
	addr := utils.GetConfig("port").(string)
	err = services.RegisterRouter(gin.Default()).Run(addr)
	if err != nil {
		panic(err)
	}
}
