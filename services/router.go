package services

import (
	"gin-blog/utils"

	"github.com/gin-gonic/gin"
)

var RouterEngine *gin.Engine

// 注册路由
func RegisterRouter(r *gin.Engine) *gin.Engine {
	RouterEngine = r
	r.NoRoute(Handle404)

	// 静态资源
	{
		assetPath := utils.GetConfig("assets_path").(string)
		r.Static("/asset", assetPath)
	}

	// 404 和 405 路由
	{

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
