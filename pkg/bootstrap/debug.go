package bootstrap

import (
	"github.com/gin-gonic/gin"
	"novels-spider/pkg/helpers"
)

func SetDebugMode() {
	debug := helpers.GetConfig("debug").(bool)
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}
