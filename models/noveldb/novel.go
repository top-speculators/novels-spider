package noveldb

import "time"

type Novel struct {
	Id                 uint
	Channel            uint   // 第三方渠道，0-未知；1-biquge.tv
	SourceId           uint   // 渠道资源 ID
	Author             string // 作者
	Name               string // 小说名
	Pic                string // 小说封面
	Intro              string // 小说简介
	MaxSequence        uint   // 最新章节的 sequence，每次更新章节都要更新此字段
	Weight             int    // 权重，值越大，越靠前展示
	Tag                int    // 0-没有标记；1-连载中；2-已完结
	Views              int    // 点击次数
	Status             int    // 状态 0-连载中 1-完结
	ChapterTableNumber int    // 分表编号
	Href               string // 源站 url
	IsOnline           int    // 0-未上线 1-已上线
	IsDeleted          int
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
