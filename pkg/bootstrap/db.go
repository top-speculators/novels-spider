package bootstrap

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"novels-spider/models/blogdb"
	"novels-spider/models/noveldb"
)

var db *gorm.DB
var novelDbs map[string]*gorm.DB

// 目前 gorm 未提供读写分离功能，需自己实现，或改用 xorm，或直接使用 mycat 等数据库中间件
func LoadDBConnections() {
	// blog sqlite3
	var err error
	db, err = blogdb.Conn()
	if err != nil {
		logrus.Error(err)
		return
	}

	// novel 库读写分离
	novelDbs, err = noveldb.Conn()
	if err != nil {
		logrus.Error(err)
		return
	}
}

func DBClose() {
	_ = db.Close()
	for _, v := range novelDbs {
		_ = v.Close()
	}
}
