package bootstrap

import (
	"github.com/beanstalkd/beanstalk"
	"github.com/cihub/seelog"
	"novels-spider/pkg/helpers"
)

var bean *beanstalk.Conn

func LoadMQConnection() {
	err := helpers.BeanConn()
	if err != nil {
		_ = seelog.Error(err)
		return
	}
	bean = helpers.GetBeanConn()
}

func MQClose() {
	_ = bean.Close()
}
