package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
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
