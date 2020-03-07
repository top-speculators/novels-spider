package helpers

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

/************************************/
/**********    文件配置相关    ********/
/************************************/

// 载入配置
var config = make(map[string]interface{})

func LoadConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err == nil {
		err = yaml.Unmarshal(data, config)
	}

	return err
}

// 获取配置值
func GetConfig(s string) interface{} {
	return config[s]
}
