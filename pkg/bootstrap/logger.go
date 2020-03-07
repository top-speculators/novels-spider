package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"novels-spider/conf"
	"os"
)

func LoadLogger() {
	// 设置日志格式为 json 格式
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// 设置输出位置
	f, err := os.Create(conf.LogCommonFile)
	if err != nil {
		logrus.Println("日志文件创建失败", err)
		return
	}

	logrus.SetOutput(f)

	// 处理 gin 的日志输出
	f1, err1 := os.Create(conf.LogCommonFileForGin)
	if err1 != nil {
		logrus.Println("日志文件创建失败", err1)
		return
	}
	gin.DefaultWriter = io.MultiWriter(f1, os.Stdout)

	// 处理 gin recover 日志
	f2, err2 := os.Create(conf.LogErrorFile)
	if err2 != nil {
		logrus.Println("日志文件创建失败", err2)
		return
	}
	gin.DefaultErrorWriter = f2
}
