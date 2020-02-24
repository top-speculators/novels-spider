package crons

import (
	"fmt"
	"gin-blog/models/noveldb"
	"github.com/PuerkitoBio/goquery"
	"github.com/cihub/seelog"
	"strconv"
	"strings"
	"sync"
	"time"
)

// TODO:不同的源站，封装成不同的结构体
// TODO:请求协程池

// 已放入 mq 等待抓取的小说列表
// 检查新小说的时候，需分别对比此变量与库中数据
// 另外，每当消费者抓取完一本时，应使用 delete 函数删除该 key
var novelsWorkingMap = make(map[string]int)
var Mu sync.Mutex

// 互斥形式得到 novelsWorkingMap
// 注意接收方一定要 mu.Unlock()
func GetNovelsWorkingMap() (nwm map[string]int, mu *sync.Mutex) {
	Mu.Lock()
	mu = &Mu
	nwm = novelsWorkingMap
	return
}

func CheckNovel() {
	// 获取源站小说列表
	_, err := GetOnlineNovels()
	if err != nil {
		_ = seelog.Error(err)
		return
	}

	// 取出目前库中所有小说
	var novels []*noveldb.Novel
	noveldb.DBs["read"].Select("name,author").Find(&novels)

	// 按书名和作者，整理成 map 格式
	novelsMap := make(map[string]int)
	for _, v := range novels {
		novelsMap[v.Name+"-"+v.Author] = 1
	}

	// 获取处理中的列表
	nwm, mu := GetNovelsWorkingMap()
	defer mu.Unlock()

	fmt.Println(nwm)
}

// 获取源站小说列表
func GetOnlineNovels() (novels map[string]string, err error) {
	startTime := time.Now()
	novels = make(map[string]string)

	paths := []string{
		"http://www.biquge.tv/xuanhuanxiaoshuo/1_1.html",
		"http://www.biquge.tv/xiuzhenxiaoshuo/2_1.html",
		"http://www.biquge.tv/dushixiaoshuo/3_1.html",
		"http://www.biquge.tv/chuanyuexiaoshuo/4_1.html",
		"http://www.biquge.tv/wangyouxiaoshuo/5_1.html",
		"http://www.biquge.tv/kehuanxiaoshuo/6_1.html",
		"http://www.biquge.tv/wanben/1_1",
	}

	// 取每分类页的分页 url
	fmt.Println("=====================> 共" + strconv.Itoa(len(paths)) + "个分类")
	for k, v := range paths {
		startTime := time.Now()
		fmt.Println("=====================> 当前第" + strconv.Itoa(k+1) + "个分类")
		l, err := GetNovelCatePageList(v)
		if err != nil {
			return nil, err
		}

		// 取每一页的小说列表
		fmt.Println("=====================> 当前分类共" + strconv.Itoa(len(l)) + "页")
		for k, v := range l {
			if k%50 == 0 && k != 0 {
				fmt.Println("=====================> 已读取 50 页")
			} else {
				if k == len(l)-1 {
					fmt.Println("=====================> 当前分类共" + strconv.Itoa(k+1) + "页读取完毕")
				}
			}
			m, err := GetNovelList(v)
			if err != nil {
				return nil, err
			}

			// 将小说添加到结果数组中
			for k, v := range m {
				novels[k] = v
			}
		}
		endTime := time.Now()                 // 结束时间
		latencyTime := endTime.Sub(startTime) // 执行时间
		fmt.Printf("=====================> 当前分类读取完毕，耗时 %13v \n", latencyTime)
		fmt.Println("=====================> 开始读取下一分类")
	}

	endTime := time.Now()                 // 结束时间
	latencyTime := endTime.Sub(startTime) // 执行时间
	fmt.Printf("任务执行完毕，共采集 %d 本小说，共耗时 %13v \n", len(novels), latencyTime)
	return
}

// 解析源站某分类所有列表页的 url
func GetNovelCatePageList(path string) (list []string, err error) {
	doc, err1 := H.GetDocumentByHttpGet(path)
	if err1 != nil {
		return nil, err1
	}

	// 总页码
	totalStr := doc.Find(".last").Text()
	total, _ := strconv.Atoi(totalStr)

	// 生成该规则下的所有 url
	for i := 1; i <= total; i++ {
		var url string
		if strings.Contains(path, ".html") {
			url = strings.Replace(path, "1.html", strconv.Itoa(i)+".html", 1)
		} else {
			url = strings.Replace(path, "_1", "_"+strconv.Itoa(i), 1)
		}
		list = append(list, url)
	}

	return
}

// 取每一页的小说列表
func GetNovelList(path string) (m map[string]string, err error) {
	doc, err1 := H.GetDocumentByHttpGet(path)
	if err1 != nil {
		return nil, err1
	}

	m = make(map[string]string)

	doc.Find(".l li").Each(func(i int, s *goquery.Selection) {
		author := s.Find(".s5").Text()
		name := s.Find(".s2 a").Text()
		href, _ := s.Find(".s2 a").Attr("href")
		author, _ = H.GBKToUTF8(author)
		name, _ = H.GBKToUTF8(name)

		m[name+"-"+author] = href
	})
	return
}
