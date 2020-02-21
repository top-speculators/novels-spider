package services

import (
	"gin-blog/interfaces"

	"github.com/gin-gonic/gin"
)

var RouterEngine *gin.Engine
var pageNum uint64
var H interfaces.Helper

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
