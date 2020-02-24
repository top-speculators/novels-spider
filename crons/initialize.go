package crons

import "gin-blog/interfaces"

var H interfaces.Helper

func SetHelper(h interfaces.Helper) {
	H = h
}
