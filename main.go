package main

import (
	"gin-blog/models"
	"gin-blog/services"
	"gin-blog/utils"

	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {

	// 日志
	{
		logger, err := seelog.LoggerFromConfigAsFile("./seelog_config.xml")
		if err != nil {
			_ = seelog.Critical("seelog 配置文件错误", err)
			return
		}

		err = seelog.ReplaceLogger(logger)
		if err != nil {
			_ = seelog.Critical("seelog 配置文件错误", err)
			return
		}

		defer seelog.Flush()
	}

	// 加载辅助函数实例
	h := utils.New()

	// 本地项目配置
	{
		// 这里的配置生命周期是 main goroutine，且全局有效
		// 当子 goroutine 修改了配置，将作用于之后开启的所有其他 goroutine ，除非 main goroutine 重启
		// 当修改了配置文件，需重启项目来使其生效
		// 注意这个配置，无法生效在其他包的 init 函数，因为那时候还未载入
		err := h.LoadConfig("./config.yaml")
		if err != nil {
			_ = seelog.Critical("本地配置加载错误", err)
			return
		}
	}

	// debug 设置
	{
		debug := h.GetConfig("debug").(bool)
		if debug {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}
	}

	// 数据库
	{
		db, _ := models.Conn(h)
		defer func() {
			_ = db.Close()
		}()
	}

	// 路由及服务
	{
		addr := h.GetConfig("port").(string)
		err := services.RegisterRouter(gin.New(), h).Run(addr)
		if err != nil {
			_ = seelog.Critical("http 服务启动错误", err)
			return
		}
	}
}
