package services

import (

	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func Test(c *gin.Context) {
	// Request the HTML page.
	res, err := http.Get("http://www.biquge.tv/xiaoshuodaquan/")
	if err != nil {
		_ = seelog.Critical("get failed", err)
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != 200 {
		_ = seelog.Critical(res.StatusCode, res.Status)
		return
	}

	// Load the HTML document
	doc, err1 := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		_ = seelog.Critical("document failed", err1)
		return
	}
	if doc == nil {
		_ = seelog.Critical("document is nil")
		return
	}

	list := make(map[string]interface{})
	// Find the review items
	doc.Find(".novellist li").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		name := s.Find("a").Text()
		href, _ := s.Find("a").Attr("href")
		hrefUTF, _ := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(href))
		nameUTF, _ := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(name))
		list[string(nameUTF)] = string(hrefUTF)
	})

	Success(c, list)
}
