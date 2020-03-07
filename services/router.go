package services

import (
	"novels-spider/services/portal/blog"
	"time"

	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

var RouterEngine *gin.Engine

// 日志中间件
// 需要通过重写 gin.Logger 中间件来将 seelog 和 gin 结合起来
func Logger(c *gin.Context) {

	// 开始时间
	startTime := time.Now()

	// 处理请求
	c.Next()

	endTime := time.Now()                 // 结束时间
	latencyTime := endTime.Sub(startTime) // 执行时间
	reqMethod := c.Request.Method         // 请求方式
	reqUri := c.Request.RequestURI        // 请求路由
	statusCode := c.Writer.Status()       // 状态码
	clientIP := c.ClientIP()              // 请求IP

	// 日志格式
	seelog.Infof("| %3d | %13v | %15s | %s | %s |",
		statusCode,
		latencyTime,
		clientIP,
		reqMethod,
		reqUri,
	)

	// TODO：完整请求记录
}

// 注册路由
func RegisterRouter(r *gin.Engine) *gin.Engine {

	// 包范围内初始化引擎
	RouterEngine = r
	r.Use(gin.Recovery())
	r.Use(Logger)
	r.Use(gin.Logger())

	r.NoRoute(blog.Handle404)

	// 前台路由 API
	{
		// 首页
		r.GET("/index", blog.Index)
		// 专辑列表
		r.GET("/categories", blog.Categories)
		// 文章列表
		r.GET("/articles/:cate", blog.Articles)
		// 文章详情
		r.GET("/article/:id", blog.Article)
		// 关于我
		r.GET("/about", blog.About)
	}

	// 后台路由
	{
		// Auth 中间件

		// 单文件上传

		// 后台首页

		// 专辑管理

		// 文章管理

		// 配置管理
	}

	return RouterEngine
}
