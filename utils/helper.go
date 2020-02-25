package utils

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/beanstalkd/beanstalk"
	"golang.org/x/text/encoding/simplifiedchinese"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type helper struct {
	config map[string]interface{}
}

func New() *helper {
	return &helper{}
}

/************************************/
/**********    爬虫函数相关    ********/
/************************************/

// 模拟 User-Agent
func (h *helper) GetRandomUserAgent() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var UserAgentList = []string{
		"Mozilla/5.0 (compatible, MSIE 10.0, Windows NT, DigExt)",
		"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, 360SE)",
		"Mozilla/4.0 (compatible, MSIE 8.0, Windows NT 6.0, Trident/4.0)",
		"Mozilla/5.0 (compatible, MSIE 9.0, Windows NT 6.1, Trident/5.0,",
		"Opera/9.80 (Windows NT 6.1, U, en) Presto/2.8.131 Version/11.11",
		"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, TencentTraveler 4.0)",
		"Mozilla/5.0 (Windows, U, Windows NT 6.1, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"Mozilla/5.0 (Macintosh, Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
		"Mozilla/5.0 (Macintosh, U, Intel Mac OS X 10_6_8, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"Mozilla/5.0 (Linux, U, Android 3.0, en-us, Xoom Build/HRI39) AppleWebKit/534.13 (KHTML, like Gecko) Version/4.0 Safari/534.13",
		"Mozilla/5.0 (iPad, U, CPU OS 4_3_3 like Mac OS X, en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5",
		"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, Trident/4.0, SE 2.X MetaSr 1.0, SE 2.X MetaSr 1.0, .NET CLR 2.0.50727, SE 2.X MetaSr 1.0)",
		"Mozilla/5.0 (iPhone, U, CPU iPhone OS 4_3_3 like Mac OS X, en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5",
		"MQQBrowser/26 Mozilla/5.0 (Linux, U, Android 2.3.7, zh-cn, MB200 Build/GRJ22, CyanogenMod-7) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
	}
	return UserAgentList[r.Intn(len(UserAgentList))]
}

// 抓取 HTML
// TODO:引入 proxy 避免爬取次数太多被拉黑
func (h *helper) GetDocumentByHttpGet(path string) (doc *goquery.Document, err error) {
	// http client
	client := &http.Client{
		Timeout: time.Second * 2,
	}

	// http request
	var req *http.Request
	req, err = http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// add User-Agent
	req.Header.Add("User-Agent", h.GetRandomUserAgent())

	res, err1 := client.Do(req)
	if err1 != nil {
		err = err1
		return
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != 200 {
		err = errors.New("源站返回错误码" + strconv.Itoa(res.StatusCode))
		return
	}

	doc, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}
	if doc == nil {
		err = errors.New("抓取 doc 为 nil")
		return
	}

	return
}

// GBK 转 UTF8
func (h *helper) GBKToUTF8(html string) (str string, err error) {
	var str1 []byte
	str1, err = simplifiedchinese.GB18030.NewDecoder().Bytes([]byte(html))
	str = string(str1)
	return
}

// UTF8 转 GBK
func (h *helper) UTF8ToGBK(html string) (str string, err error) {
	var str1 []byte
	str1, err = simplifiedchinese.GB18030.NewEncoder().Bytes([]byte(html))
	str = string(str1)
	return
}

/************************************/
/**********    文件配置相关    ********/
/************************************/

// 载入配置
func (h *helper) LoadConfig(path string) error {
	h.config = make(map[string]interface{})

	data, err := ioutil.ReadFile(path)
	if err == nil {
		err = yaml.Unmarshal(data, h.config)
	}

	return err
}

// 获取配置值
func (h *helper) GetConfig(s string) interface{} {
	return h.config[s]
}

/************************************/
/**********    消息队列相关    ********/
/************************************/
var beanstalkdConn *beanstalk.Conn

// 建立连接
func (h *helper) BeanstalkdConn() error {
	c, err := beanstalk.Dial("tcp", h.GetConfig("beanstalkd_dsn").(string))
	if err != nil {
		return err
	}
	beanstalkdConn = c
	return nil
}

// 获取连接对象
func (h *helper) GetBeanstalkdConn() *beanstalk.Conn {
	return beanstalkdConn
}
