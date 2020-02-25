package interfaces

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/beanstalkd/beanstalk"
)

// 封装 Helper
// 所有用到辅助函数的地方都依赖于此接口，而不依赖于具体的包，如 utils
// 这样底层实现的包随时可替换

type Helper interface {
	// 载入配置
	// path 指配置文件路径，只支持 yaml 文件
	LoadConfig(path string) error

	// 获取配置值
	// s 指配置名
	GetConfig(s string) interface{}

	// 模拟 User-Agent
	GetRandomUserAgent() string
	// 抓取网页
	// path 是 url
	GetDocumentByHttpGet(path string) (doc *goquery.Document, err error)
	// GBK 转 UTF8
	GBKToUTF8(html string) (str string, err error)
	// UTF8 转 GBK
	UTF8ToGBK(html string) (str string, err error)

	// beanstalkd 相关
	BeanstalkdConn() error
	GetBeanstalkdConn() *beanstalk.Conn
}
