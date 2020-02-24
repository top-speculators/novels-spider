package services

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

func Test(c *gin.Context) {

	doc, err := H.GetDocumentByHttpGet("http://www.biquge.tv/xiaoshuodaquan/")
	if err != nil {
		_ = seelog.Critical(err)
		return
	}
	list := make(map[string]interface{})

	// 匹配内容
	doc.Find(".novellist li").Each(func(i int, s *goquery.Selection) {
		novelName := s.Find("a").Text()
		novelHref, _ := s.Find("a").Attr("href")
		novelName, _ = H.GBKToUTF8(novelName)
		novelHref, _ = H.GBKToUTF8(novelHref)
		list[novelName] = novelHref
	})

	Success(c, gin.H{
		"content": list,
	})
}
