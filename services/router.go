package services

import (
	"gin-blog/interfaces"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"time"
)

var RouterEngine *gin.Engine
var pageNum uint64
var H interfaces.Helper

// 日志中间件
// 需要通过重写 gin.Logger 中间件来将 seelog 和 gin 结合起来
func Logger(c *gin.Context) {

	// 开始时间
	startTime := time.Now()

	// 处理请求
	c.Next()

	endTime := time.Now() // 结束时间
	latencyTime := endTime.Sub(startTime) // 执行时间
	reqMethod := c.Request.Method // 请求方式
	reqUri := c.Request.RequestURI // 请求路由
	statusCode := c.Writer.Status() 	// 状态码
	clientIP := c.ClientIP() // 请求IP

	// 日志格式
	seelog.Infof("| %3d | %13v | %15s | %s | %s |",
		statusCode,
		latencyTime,
		clientIP,
		reqMethod,
		reqUri,
	)
}

// 注册路由
func RegisterRouter(r *gin.Engine, h interfaces.Helper) *gin.Engine {
	H = h

	numStr := H.GetConfig("page_num").(int)
	pageNum = uint64(numStr)

	RouterEngine = r
	r.Use(gin.Recovery())
	r.Use(Logger)

	r.NoRoute(Handle404)

	// 静态资源
	{
		assetPath := H.GetConfig("assets_path").(string)
		r.Static("/asset", assetPath)
	}

	// 前台路由
	{
		// 首页
		r.GET("/", Index)
		// 专辑列表
		r.GET("/categories", Categories)
		// 文章列表
		r.GET("/articles/:cate", Articles)
		// 文章详情
		r.GET("/article/:id", Article)
		// 关于我
		r.GET("/about", About)
	}

	// 后台路由
	{
		// BaseAuth 中间件

		// 单文件上传

		// 后台首页

		// 专辑管理

		// 文章管理

		// 配置管理
	}

	return RouterEngine
}
