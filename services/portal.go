package services

import (
	"errors"
	"strconv"

	"gin-blog/models"

	"github.com/gin-gonic/gin"
)

// 前台各页面接口通用数据
func commonData() (links []*models.Link, categories []*models.Category, siteMap map[string]interface{}) {

	models.DB.Order("visits desc").Limit(6).Find(&categories)
	models.DB.Find(&links)

	siteMap = make(map[string]interface{})
	siteMap["siteTitle"] = models.GetSiteConfig("site_title")
	siteMap["siteCopyRight"] = models.GetSiteConfig("site_copyRight")
	siteMap["siteCountCode"] = models.GetSiteConfig("site_count_code")

	return
}

// 页码处理
func getPage(c *gin.Context) uint64 {

	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.ParseUint(pageStr, 0, 0)
	if err != nil {
		page = 1
	}
	return page
}

// 404 处理
func Handle404(c *gin.Context) {
	Failed(c, 404, errors.New("资源不存在"))
	c.Abort() // 避免后面 handlers 被调用
}

// 首页 html
func IndexHtml(c *gin.Context) {
	path := H.GetConfig("portal_html").(string)
	RouterEngine.LoadHTMLFiles(path)
	c.HTML(200, "portal.html", gin.H{})
}

// 首页
func Index(c *gin.Context) {
	articleModel := &models.Article{}
	articles := articleModel.GetList(1, pageNum, "")
	links, categories, siteMap := commonData()

	Success(c, gin.H{
		"articles":   articles,
		"links":      links,
		"categories": categories,
		"siteMap":    siteMap,
	})
}

// 分类列表页
func Categories(c *gin.Context) {
	categoryModel := &models.Category{}
	categories := categoryModel.GetList(getPage(c), pageNum, "")
	pageCount := categoryModel.PageCount(pageNum, "")

	links, _, siteMap := commonData()

	Success(c, gin.H{
		"links":      links,
		"categories": categories,
		"siteMap":    siteMap,
		"pageCount":  pageCount,
	})
}

// 文章列表页
func Articles(c *gin.Context) {

	cate := c.Param("cate")
	categoryModel := &models.Category{}
	category := categoryModel.First("name = ?", cate)
	if category.ID == 0 {
		Handle404(c)
		return // 避免此 handler 后面的代码被调用
	}

	articleModel := &models.Article{}
	articles := articleModel.GetList(getPage(c), pageNum, "category_id = ?", category.ID)
	pageCount := articleModel.PageCount(pageNum, "")

	links, categories, siteMap := commonData()

	Success(c, gin.H{
		"links":      links,
		"categories": categories,
		"siteMap":    siteMap,
		"pageCount":  pageCount,
		"articles":   articles,
		"category":   category,
	})
}

// 文章页
func Article(c *gin.Context) {

	id := c.Param("id")
	articleModel := &models.Article{}
	article := articleModel.First("id = ?", id)
	categoryModel := &models.Category{}
	category := categoryModel.First("id = ?", article.CategoryId)

	links, categories, siteMap := commonData()

	Success(c, gin.H{
		"links":      links,
		"categories": categories,
		"siteMap":    siteMap,
		"category":   category,
		"article":    article,
	})
}

// 关于我
func About(c *gin.Context) {

	links, categories, siteMap := commonData()

	Success(c, gin.H{
		"links":      links,
		"categories": categories,
		"siteMap":    siteMap,
	})
}
