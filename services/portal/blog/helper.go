package blog

import (
	"github.com/gin-gonic/gin"
	"novels-spider/models/blogdb"
)

var pageNum uint64 = 10 // 一页默认条数

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

// 前台各页面接口通用数据
func commonData() (links []*blogdb.Link, categories []*blogdb.Category, siteMap map[string]interface{}) {

	blogdb.DB.Order("visits desc").Limit(6).Find(&categories)
	blogdb.DB.Find(&links)

	siteMap = make(map[string]interface{})
	siteMap["siteTitle"] = blogdb.GetSiteConfig("site_title")
	siteMap["siteCopyRight"] = blogdb.GetSiteConfig("site_copyRight")
	siteMap["siteCountCode"] = blogdb.GetSiteConfig("site_count_code")

	return
}