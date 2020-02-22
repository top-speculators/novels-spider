package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
)

type helper struct {
	config map[string]interface{}
}

func New() *helper {
	return &helper{}
}

/************************************/
/**********    基本辅助函数    ********/
/************************************/

// 设置日志处理
// t = true 	: 同时开启控制台日志和文件日志
// t = false 	: 只开启文件日志
// io.MultiWriter 操作文件时会占用该文件，main.go 运行期间无法删除和重命名日志文件，但可以修改日志文件的内容，会引起乱码
func (h *helper) SetLogConfig(t bool, path string) {

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

// 载入配置
func (h *helper) LoadConfig(path string) error {
	h.config = make(map[string]interface{})

	data, err := ioutil.ReadFile(path)
	if err == nil {
		err = yaml.Unmarshal(data, h.config)
	}

	return err
}

// 获取配置值
func (h *helper) GetConfig(s string) interface{} {
	return h.config[s]
}
