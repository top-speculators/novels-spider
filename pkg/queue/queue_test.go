package queue

// ----------------------------------------------------------------------
// queue 测试
// ----------------------------------------------------------------------
// Link  : http://www.hlzblog.top/
// GITHUB: https://github.com/HaleyLeoZhang
// ----------------------------------------------------------------------
// tesing包文档 https://golang.google.cn/pkg/testing/
// ----------------------------------------------------------------------

import (
	"novels-spider/pkg/logging"
	"testing"
)

func TestMain(m *testing.M) {
	logging.Setup()
	m.Run()
}

func TestCore(t *testing.T) {
	t.Run("Push", Push)
	t.Run("Pull", Pull)
}

func Push(t *testing.T) {
	q := GetNovelQueue()
	payload := "获取一本小说"
	q.SetPayload([]byte(payload))

	err := q.Push()
	if err != nil {
		t.Fatalf("生产消息失败: %s \n", err)
	} else {
		logging.Debug("Push:", string(payload))
	}
}

var testFlagClose chan int
var counter int

func Pull(t *testing.T) {
	testFlagClose = make(chan int, 1)

	q := GetNovelQueue()
	go func() {
		err := q.Pull(callPull)
		if err != nil {
			t.Fatalf("消费消息失败: %v ", err)
		}
	}()
	<-testFlagClose
	q.Close() // 记得单元测试关闭连接
}

func callPull(payload []byte) error {
	logging.Debug("Pull:", string(payload))
	testFlagClose <- 1
	return nil
}
