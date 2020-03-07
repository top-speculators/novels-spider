package bootstrap

import (
	"github.com/gin-gonic/gin"
	"novels-spider/conf"
)

func SetDebugMode() {
	if conf.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}