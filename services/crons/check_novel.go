package crons

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"novels-spider/pkg/helpers"
	"strconv"
	"strings"
	"sync"
	"time"

	"novels-spider/models/noveldb"

	"github.com/PuerkitoBio/goquery"
)

// TODO:不同的源站，封装成不同的结构体

// novelsWorkingMap 代表已放入 mq 等待抓取的小说列表
// 检查新小说的时候，需分别对比此变量与库中数据
// 另外，每当消费者抓取完一本时，应使用 delete 函数删除该 key
var novelsWorkingMap = make(map[string]uint64)
var Mu sync.Mutex
var novelOnlineChannel = 1 // novel 在线渠道，1 为 biquge.tv

// 互斥形式得到 novelsWorkingMap
// 注意接收方一定要 mu.Unlock()
func GetNovelsWorkingMap() (nwm map[string]uint64, mu *sync.Mutex) {
	Mu.Lock()
	mu = &Mu
	nwm = novelsWorkingMap
	return
}

func CheckNovel() {
	// 获取源站小说列表
	// 先爬取数据，避免占用 novelsWorkingMap 太久，影响到其他 goroutine
	onlineNovels, err := GetOnlineNovels()
	if err != nil {
		logrus.Error(err)
		return
	}

	// 取出目前库中所有小说
	var novels []*noveldb.Novel
	noveldb.DBs["read"].Select("name,author").Find(&novels)

	// 按书名和作者，整理成 map 格式
	// 已在库的，用不到 url，所以 value 定义为 int 类型
	novelsMap := make(map[string]uint64)
	for _, v := range novels {
		// name 和 author 才能唯一决定一部小说
		novelsMap[v.Name+"-"+v.Author] = 1
	}

	// 获取处理中的列表
	nwm, mu := GetNovelsWorkingMap()
	defer mu.Unlock() // 在对比现有 novel 和 online novel 的整个过程，nwm 都会被阻塞

	for k, v := range onlineNovels {
		_, ok := nwm[k]
		_, ok2 := novelsMap[k]
		if !ok && !ok2 {
			// 如果两个列表里都不存在

			// 往消息队列里生产 Job
			newNovelTube := helpers.GetBeanTube(NewNovelTube)

			// name-author:href:channel
			job := k + ":" + v + ":" + string(novelOnlineChannel)
			id, err := newNovelTube.Put([]byte(job), 1, 0, 120*time.Second) // 120 秒后触发 TTR
			if err != nil {
				// 写入 mq 失败，继续直接进入下一层循环
				logrus.Error(err, job)
				continue
			}
			novelsMap[k] = id // 往 nwm 中写入数据，等消费者执行结束后，再删除
		}
	}
}

// 获取源站小说列表
func GetOnlineNovels() (novels map[string]string, err error) {
	startTime := time.Now()
	novels = make(map[string]string)

	// 源站分类 url 数组切片
	paths := []string{
		"http://www.biquge.tv/xuanhuanxiaoshuo/1_1.html",
		"http://www.biquge.tv/xiuzhenxiaoshuo/2_1.html",
		"http://www.biquge.tv/dushixiaoshuo/3_1.html",
		"http://www.biquge.tv/chuanyuexiaoshuo/4_1.html",
		"http://www.biquge.tv/wangyouxiaoshuo/5_1.html",
		"http://www.biquge.tv/kehuanxiaoshuo/6_1.html",
		"http://www.biquge.tv/wanben/1_1",
	}

	fmt.Println("=====================> 共" + strconv.Itoa(len(paths)) + "个分类")

	// 此通道用于并发获取小说所开启时子协程到父协程的数据回传
	// 这里的父子是逻辑关系，物理上，go 只有主协程和子协程
	// 指定 1000 的容量，避免因为接收者定义在 “开启协程” 循环体的下面
	// 而导致所有开启的协程都会因为没有接收者而阻塞的情况
	// 若没有指定容量，则开启的 N 个协程，都会在等到接收者出现时，再统一释放
	// 特别是在 “开启协程” 循环体内还存在 time.Sleep 的情况下，更应该注意
	// 1000 的容量，则可以提前释放 1000 个 goroutine
	// 之所以容量定 1000 ，考虑到所有分类下的所有分页，全部加起来，不超过 1000（目前源站是 973 页）
	ch := make(chan map[string]string, 1000)
	defer close(ch)

	// 开启的总协程数，用于判断 ch 的关闭时机
	var chsCount int

	// 按小说分类遍历
	for k, v := range paths {
		startTime := time.Now()
		fmt.Println("=====================> 当前第" + strconv.Itoa(k+1) + "个分类")

		// 组织该分类下的所有分页 url
		// 此方法为同步阻塞式请求
		l, err := GetNovelCatePages(v)
		if err != nil {
			return nil, err
		}

		fmt.Println("=====================> 当前分类共" + strconv.Itoa(len(l)) + "页")

		// 遍历该分类下的所有分页，解析每一张 dom ，从中得到该页面上所出现的小说列表
		for _, v := range l {
			// 避免源站 timeout 限制 1 秒 15 并发请求
			if chsCount%15 == 1 && chsCount != 1 {
				fmt.Println("=====================> 已开启 15 个协程，先暂停 1 秒")
				time.Sleep(time.Second)
			}

			// 解析该页小说列表
			go GetNovelsByCatePage(v, ch)
			chsCount++
		}

		endTime := time.Now()                 // 结束时间
		latencyTime := endTime.Sub(startTime) // 执行时间
		fmt.Printf("=====================> 当前分类读取完毕，耗时 %13v \n", latencyTime)
	}

	// 处理回传的小说列表
	for {
		if chsCount > 0 {
			m := <-ch
			chsCount--
			for k, v := range m {
				novels[k] = v
			}
		} else {
			break
		}
	}

	endTime := time.Now()                 // 结束时间
	latencyTime := endTime.Sub(startTime) // 执行时间
	fmt.Printf("任务执行完毕，共采集 %d 本小说，共耗时 %13v \n", len(novels), latencyTime)
	return
}

// 解析源站某分类下的所有分页 url
func GetNovelCatePages(path string) (list []string, err error) {
	doc, err1 := helpers.GetDocumentByHttpGet(path)
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

// 解析某分类页下的小说列表
func GetNovelsByCatePage(path string, ch chan<- map[string]string) {
	doc, err1 := helpers.GetDocumentByHttpGet(path)
	if err1 != nil {
		logrus.Error(err1)
		return
	}

	m := make(map[string]string)

	doc.Find(".l li").Each(func(i int, s *goquery.Selection) {
		author := s.Find(".s5").Text()
		name := s.Find(".s2 a").Text()
		href, _ := s.Find(".s2 a").Attr("href")
		author, _ = helpers.GBKToUTF8(author)
		name, _ = helpers.GBKToUTF8(name)

		m[name+"-"+author] = href
	})

	ch <- m

	// 可以在此处添加打印，来观察给 ch 指定容量和不给 ch 指定容量的不同现象
}
