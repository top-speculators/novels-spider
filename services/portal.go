package services

import (
	"gin-blog/utils"

	"strconv"

	"gin-blog/models"
	"github.com/gin-gonic/gin"
)

// 默认分页
var pageNum = utils.GetConfig("page_num").(uint64)

// 前台页面通用数据
func commonData() (links []*models.Link, categories []*models.Category, siteMap map[string]interface{}) {

	models.DB.Order("visits desc").Limit(6).Find(&categories)
	models.DB.Find(&links)

	siteMap = make(map[string]interface{})
	siteMap["siteTitle"] = models.GetSiteConfig("site_title")
	siteMap["siteDescription"] = models.GetSiteConfig("site_description")
	siteMap["siteKeyword"] = models.GetSiteConfig("site_keyword")
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
	c.String(404, "404 not found")
	c.Abort() // 避免后面 handlers 被调用
}

// 首页
func Index(c *gin.Context) {
	articleModel := &models.Article{}
	articles := articleModel.GetList(1, pageNum, "")
	links, categories, siteMap := commonData()

	c.JSON(200, gin.H{
		"articles":   articles,
		"links":      links,
		"categories": categories,
		"siteMap":    siteMap,
	})

	//indexTmpl := utils.GetConfig("index_html").(string)
	//RouterEngine.LoadHTMLFiles(indexTmpl)
	//c.HTML(http.StatusOK, indexTmpl, gin.H{
	//	"articles":   articles,
	//	"links":      links,
	//	"categories": categories,
	//	"siteMap":    siteMap,
	//})
}

// 分类列表页
func Categories(c *gin.Context) {
	categoryModel := &models.Category{}
	categories := categoryModel.GetList(getPage(c), pageNum, "")
	pageCount := categoryModel.PageCount(pageNum, "")

	links, _, siteMap := commonData()

	c.JSON(200, gin.H{
		"links":      links,
		"categories": categories,
		"siteMap":    siteMap,
		"pageCount":  pageCount,
	})

	//indexTmpl := utils.GetConfig("categories_html").(string)
	//RouterEngine.LoadHTMLFiles(indexTmpl)
	//c.HTML(http.StatusOK, indexTmpl, gin.H{
	//	"links":      links,
	//	"categories": categories,
	//	"siteMap":    siteMap,
	//	"pageCount":  pageCount,
	//})
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

	c.JSON(200, gin.H{
		"links":      links,
		"categories": categories,
		"siteMap":    siteMap,
		"pageCount":  pageCount,
		"articles":   articles,
		"category":   category,
	})

	//indexTmpl := utils.GetConfig("articles_html").(string)
	//RouterEngine.LoadHTMLFiles(indexTmpl)
	//c.HTML(http.StatusOK, indexTmpl, gin.H{
	//	"links":      links,
	//	"categories": categories,
	//	"siteMap":    siteMap,
	//	"pageCount":  pageCount,
	//	"articles":   articles,
	//	"category":   category,
	//})
}

// 文章页
func Article(c *gin.Context) {

	id := c.Param("id")
	articleModel := &models.Article{}
	article := articleModel.First("id = ?", id)
	categoryModel := &models.Category{}
	category := categoryModel.First("id = ?", article.CategoryId)

	links, categories, siteMap := commonData()

	c.JSON(200, gin.H{
		"links":      links,
		"categories": categories,
		"siteMap":    siteMap,
		"category":   category,
		"article":    article,
	})

	//indexTmpl := utils.GetConfig("article_html").(string)
	//RouterEngine.LoadHTMLFiles(indexTmpl)
	//c.HTML(http.StatusOK, indexTmpl, gin.H{
	//	"links":      links,
	//	"categories": categories,
	//	"siteMap":    siteMap,
	//	"category":   category,
	//	"article":    article,
	//})
}

// 关于我
func About(c *gin.Context) {

	links, categories, siteMap := commonData()

	c.JSON(200, gin.H{
		"links":      links,
		"categories": categories,
		"siteMap":    siteMap,
	})

	//indexTmpl := utils.GetConfig("about)html").(string)
	//RouterEngine.LoadHTMLFiles(indexTmpl)
	//c.HTML(http.StatusOK, indexTmpl, gin.H{
	//	"links":      links,
	//	"categories": categories,
	//	"siteMap":    siteMap,
	//})
}
