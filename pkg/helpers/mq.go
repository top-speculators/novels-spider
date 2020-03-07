package helpers

import (
	"github.com/beanstalkd/beanstalk"
	"novels-spider/conf"
)

var beanConn *beanstalk.Conn

// 建立连接
func BeanConn() error {
	c, err := beanstalk.Dial("tcp", conf.Beanstalkd)
	if err != nil {
		return err
	}
	beanConn = c
	return nil
}

// 获取连接对象
func GetBeanConn() *beanstalk.Conn {
	return beanConn
}

// 获取一个 Tube
func GetBeanTube(name string) *beanstalk.Tube {
	return &beanstalk.Tube{Conn: beanConn, Name: name}
}
