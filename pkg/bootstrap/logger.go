package bootstrap

import "github.com/cihub/seelog"

// 加载日志
func LoadLogger() {
	logger, err := seelog.LoggerFromConfigAsFile("./seelog_config.xml")
	if err != nil {
		_ = seelog.Critical("seelog 配置文件错误", err)
		return
	}

	err = seelog.ReplaceLogger(logger)
	if err != nil {
		_ = seelog.Critical("seelog 配置文件错误", err)
		return
	}
}

// 刷新日志
func LoggerFlush() {
	seelog.Flush()
}
