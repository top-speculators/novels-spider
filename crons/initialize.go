package crons

import "gin-blog/interfaces"

var NewNovelTube = "newNovel"
var ChapterUpdaterTube = "chapterUpdater"

var H interfaces.Helper

func SetHelper(h interfaces.Helper) {
	H = h
}
