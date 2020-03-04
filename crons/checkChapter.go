package crons

import "gin-blog/models/noveldb"

func CheckChapter() {

	// 拉取所有库中的 novel
	var novels []*noveldb.Novel
	noveldb.DBs["read"].Select("name,author").Find(&novels)

	// 循环检查 novel 是否有更新，开启 goroutine 检查，每 15 个协程，暂停 2 秒
	count := len(novels)
	for _, v := range novels {

		count--
	}

	// 这里不检查 mq 中的 job
	// 因为考虑到只有在一个 cron 周期消费者还未消费完 job 的情况下，才会重合，重合率很低
	// 所以只在消费者做数据入库时保证幂等性即可
}
