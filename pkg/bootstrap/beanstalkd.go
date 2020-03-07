package bootstrap

import (
	"github.com/beanstalkd/beanstalk"
	"github.com/sirupsen/logrus"
	"novels-spider/pkg/helpers"
)

var bean *beanstalk.Conn

func LoadMQConnection() {
	err := helpers.BeanConn()
	if err != nil {
		logrus.Error(err)
		return
	}
	bean = helpers.GetBeanConn()
}

func MQClose() {
	_ = bean.Close()
}
