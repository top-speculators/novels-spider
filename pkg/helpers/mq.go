package helpers

import (
	"github.com/beanstalkd/beanstalk"
)

/************************************/
/**********    消息队列相关    ********/
/************************************/
var beanConn *beanstalk.Conn

// 建立连接
func BeanConn() error {
	c, err := beanstalk.Dial("tcp", GetConfig("beanstalkd_dsn").(string))
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
