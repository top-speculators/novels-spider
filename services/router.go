package services

import (
	"time"

	"gin-blog/interfaces"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
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

// 成功响应
func Success(c *gin.Context, data gin.H) {
	res := map[string]interface{}{
		"code":    "1",
		"message": "ok",
	}
	res["data"] = data

	c.JSON(200, res)
}

// 失败响应
func Failed(c *gin.Context, code int, err error) {
	res := map[string]interface{}{
		"code":    "-1",
		"message": err.Error(),
		"data":    make(gin.H),
	}
	c.JSON(code, res)
}

// 注册路由
func RegisterRouter(r *gin.Engine, h interfaces.Helper) *gin.Engine {
	// 包范围内初始化辅助函数包
	H = h

	// 包范围内初始化每页条数
	numStr := H.GetConfig("page_num").(int)
	pageNum = uint64(numStr)

	// 包范围内初始化引擎
	RouterEngine = r
	r.Use(gin.Recovery())
	r.Use(Logger)
	r.Use(gin.Logger())

	r.NoRoute(Handle404)

	// 静态资源服务
	{
		r.Static("/static", H.GetConfig("static_path").(string))
	}

	r.GET("/", IndexHtml)
	//r.GET("/admin", AdminHtml)

	// 前台路由 API
	{
		// 首页
		r.GET("/index", Index)
		// 专辑列表
		r.GET("/categories", Categories)
		// 文章列表
		r.GET("/articles/:cate", Articles)
		// 文章详情
		r.GET("/article/:id", Article)
		// 关于我
		r.GET("/about", About)
	}

	// 爬虫
	{
		r.GET("/spider/test",Test)
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
