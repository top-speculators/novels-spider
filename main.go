package main

import (
	"gin-blog/models/blogdb"
	"gin-blog/models/noveldb"
	"gin-blog/services"
	"gin-blog/services/crons"
	"gin-blog/utils"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/robfig/cron/v3"
)

func main() {

	// 加载辅助函数实例
	h := utils.New()

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
	// 目前 gorm 未提供读写分离功能，需自己实现，或改用 xorm，或直接使用 mycat 等数据库中间件
	{
		// blog sqlite3
		db, err1 := blogdb.Conn(h)
		defer func() {
			_ = db.Close()
		}()
		if err1 != nil {
			_ = seelog.Error(err1)
			return
		}

		// novel 库读写分离
		novelDbs, err2 := noveldb.Conn(h)
		defer func() {
			for _, v := range novelDbs {
				_ = v.Close()
			}
		}()
		if err2 != nil {
			_ = seelog.Error(err2)
			return
		}
	}

	// beanstalkd 连接
	{
		err := h.BeanConn()
		if err != nil {
			_ = seelog.Error(err)
			return
		}
		b := h.GetBeanConn()
		defer func() {
			_ = b.Close()
		}()
	}

	// cron
	go func() {
		c := cron.New()
		crons.SetHelper(h)
		cronId1, _ := c.AddFunc("0 1 /5 * *", crons.CheckNovel)  // 每 5 天凌晨 1 点，检查是否有新小说
		cronId2, _ := c.AddFunc("0 6 * * *", crons.CheckChapter) // 每天凌晨 6 点，检查所有小说章节是否有更新
		seelog.Infof("计时任务已开启，%d d%", cronId1, cronId2)
		c.Start()
		select {}
	}()

	// check novel customer
	go func() {

	}()

	// check chapter customer
	go func() {

	}()

	// web 服务
	addr := h.GetConfig("port").(string)
	err := services.RegisterRouter(gin.New(), h).Run(addr)
	if err != nil {
		_ = seelog.Critical("http 服务启动错误", err)
		return
	}
}
