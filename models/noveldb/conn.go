package noveldb

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"novels-spider/pkg/helpers"
	"os"
)

var DBs map[string]*gorm.DB

func Conn() (dbs map[string]*gorm.DB, err error) {
	dbs = make(map[string]*gorm.DB)

	// 建立连接
	dsnRead := helpers.GetConfig("mysql_dsn_read")
	dbs["read"], err = gorm.Open("mysql", dsnRead)
	dsnWrite := helpers.GetConfig("mysql_dsn_write")
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
	logFile := helpers.GetConfig("novel_db_log_file_" + key).(string)
	f, err := os.Create(logFile)
	if err != nil {
		fmt.Printf("get form err: %s", err.Error())
	}

	db.SetLogger(log.New(f, "\r\n", 0))
}
