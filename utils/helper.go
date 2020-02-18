package utils

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
)

/************************************/
/**********    基本辅助函数    ********/
/************************************/

// 设置日志处理
// t = true 	: 同时开启控制台日志和文件日志
// t = false 	: 只开启文件日志
// io.MultiWriter 操作文件时会占用该文件，main.go 运行期间无法删除和重命名日志文件，但可以修改日志文件的内容，会引起乱码
func SetLogConfig(t bool, path string) {

	f, err := os.Create(path)
	if err != nil {
		fmt.Printf("get form err: %s", err.Error())
	}

	if t {
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	} else {
		gin.DefaultWriter = io.MultiWriter(f)
	}

}

/************************************/
/**********    文件配置相关    ********/
/************************************/

var config map[string]interface{}

// 载入配置
func LoadConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err == nil {
		err = yaml.Unmarshal(data, &config)
	}
	return err
}

// 获取配置值
func GetConfig(key string) interface{} {
	return config[key]
}