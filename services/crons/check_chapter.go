package crons

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"novels-spider/models/noveldb"
	"novels-spider/pkg/helpers"
	"time"
)

// 检查所有库存小说
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

// 检查单本小说是否更新
// 若有更新，则将小说写入 mq
func CheckOnlineChapter(v *noveldb.Novel) {
	// 爬取
	doc, err1 := helpers.GetDocumentByHttpGet(v.Href)
	if err1 != nil {
		logrus.Error(err1)
		return
	}

	// 对比
	onlineSum := doc.Find(".box_con .list a").Size()
	var count int
	// 使用 idx-channel-source_id-sequence 索引
	noveldb.DBs["read"].Select("id").Where("channel = ?", 1).Where("source_id = ?", v.SourceId).Count(&count)
	if (onlineSum - 9) > count {
		// 有更新，生成 Job
		chapterUpdaterTube := helpers.GetBeanTube(ChapterUpdaterTube)
		// href:分表号:渠道号:资源ID:生成任务时的已有章节数
		job := v.Href + ":" + string(v.ChapterTableNumber) + ":" + string(v.Channel) + ":" + string(v.SourceId) + ":" + string(count)
		// 优先级，越早的 job 越先消费，使 job 呈队列状
		// 避免一个 cron 周期间隔后，消费者还未消费完该次 cron 所生产的所有 job，从而导致出现漏抓章节问题
		pri := 999999999999999 - time.Now().Second()
		_, err := chapterUpdaterTube.Put([]byte(job), uint32(pri), 0, 120*time.Second) // 120 秒后触发 TTR
		if err != nil {
			logrus.Error(err, job)
		}
	}

	return
}
