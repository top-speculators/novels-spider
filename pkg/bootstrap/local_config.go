package bootstrap

import (
	"github.com/cihub/seelog"
	"novels-spider/pkg/helpers"
)

// 这里的配置生命周期是 main goroutine，且全局有效
// 当子 goroutine 修改了配置，将作用于之后开启的所有其他 goroutine ，除非 main goroutine 重启
// 当修改了配置文件，需重启项目来使其生效
// 注意这个配置，无法生效在其他包的 init 函数，因为那时候还未载入
func LoadLocalConfig() {
	err := helpers.LoadConfig("./config.yaml")
	if err != nil {
		_ = seelog.Critical("本地配置加载错误", err)
		return
	}
}
