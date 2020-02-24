package noveldb

import (
	"fmt"
	"gin-blog/interfaces"
	"github.com/jinzhu/gorm"
	"log"
	"os"
)

var DBs map[string]*gorm.DB
var H interfaces.Helper

func Conn(h interfaces.Helper) (dbs map[string]*gorm.DB, err error) {
	H = h
	dbs = make(map[string]*gorm.DB)

	// 建立连接
	dsnRead := H.GetConfig("mysql_dsn_read")
	dbs["read"], err = gorm.Open("mysql", dsnRead)
	dsnWrite := H.GetConfig("mysql_dsn_write")
	dbs["write"], err = gorm.Open("mysql", dsnWrite)
	if err != nil {
		return nil, err
	}

	// 连接成功
	DBs = dbs
	for key, db := range DBs {
		db.LogMode(true)
		SetGormLogger(db, key)
	}

	return
}

// 日志输出到文件
func SetGormLogger(db *gorm.DB, key string) {
	logFile := H.GetConfig("novel_db_log_file_" + key).(string)
	f, err := os.Create(logFile)
	if err != nil {
		fmt.Printf("get form err: %s", err.Error())
	}

	db.SetLogger(log.New(f, "\r\n", 0))
}
