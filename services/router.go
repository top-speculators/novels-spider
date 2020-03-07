package services

import (
	"github.com/gin-gonic/gin"
	"novels-spider/services/portal/blog"
)

var RouterEngine *gin.Engine

// 注册路由
func RegisterRouter(r *gin.Engine) *gin.Engine {

	// 包范围内初始化引擎
	RouterEngine = r
	r.Use(gin.Recovery())
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
