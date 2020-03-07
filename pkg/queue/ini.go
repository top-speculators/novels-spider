package queue

// ----------------------------------------------------------------------
// 初始化包
// ----------------------------------------------------------------------
// Link  : http://www.hlzblog.top/
// GITHUB: https://github.com/HaleyLeoZhang
// ----------------------------------------------------------------------

import (
	"novels-spider/conf"
	"novels-spider/pkg/e"
	"sync"
)

/**
 * 包内全局变量
 */
type oneInstacne struct {
	// 单例连接
	Conn interface{}
	// 因为多协程共用一个tcp链接,防止并发交错错写入
	// - 但一个连接能建立多个通道
	Lock sync.Mutex
}

var one oneInstacne

/**
 * 简单工厂
 *
 * @return queue.Queue
 */
func GetNovelQueue() Queue {
	switch conf.QUEUE_DRIVER {
	case "amqp":
		return &AMQP{
			Exchange: e.AMQP_NOVEL_EXCHANGE,
			Queue:    e.AMQP_NOVEL_QUEUE,
		}
	case "kafka":
		return &Kafka{}
	default:
		panic("驱动配置错误")
	}
}
