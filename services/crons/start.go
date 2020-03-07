package crons

import (
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func StartListen(){
	c := cron.New()
	cronId1, _ := c.AddFunc("0 1 /5 * *", CheckNovel)  // 每 5 天凌晨 1 点，检查是否有新小说
	cronId2, _ := c.AddFunc("0 6 * * *", CheckChapter) // 每天凌晨 6 点，检查所有小说章节是否有更新
	logrus.Infof("计时任务已开启，%d %d", cronId1, cronId2)
	c.Start()
	select {} // 无需关闭，main 结束后，此协程自动释放
}
