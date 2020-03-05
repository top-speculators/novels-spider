package crons

import (
	"fmt"
	"gin-blog/models/noveldb"
	"time"
)

func CheckChapter() {

	// 拉取所有库中的 novel
	var novels []*noveldb.Novel
	noveldb.DBs["read"].Select("name,author").Find(&novels)

	fmt.Println("=====================> 开始检查 novels 的更新情况，当前库存共 " + string(len(novels)) + " 本")
	// 循环检查 novel 是否有更新
	for k, v := range novels {
		if (k+1)%15 == 1 && (k+1) != 1 {
			fmt.Println("=====================> 已开启 15 个协程，先暂停 2 秒")
			time.Sleep(2 * time.Second)
		}

		go CheckOnlineChapter(v)
	}

	// 这里不检查 mq 中的 job
	// 因为考虑到只有在一个 cron 周期消费者还未消费完 job 的情况下，才会重合，重合率很低
	// 所以只在消费者做数据入库时保证幂等性即可
}

func CheckOnlineChapter(v *noveldb.Novel) {

}
