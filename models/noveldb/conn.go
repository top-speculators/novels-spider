package noveldb

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"novels-spider/conf"
	"os"
)

var DBs map[string]*gorm.DB

func Conn() (dbs map[string]*gorm.DB, err error) {
	dbs = make(map[string]*gorm.DB)

	// 建立连接
	dsnRead := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		conf.MysqlAvatarRead["username"],
		conf.MysqlAvatarRead["password"],
		conf.MysqlAvatarRead["ip"],
		conf.MysqlAvatarRead["port"],
		conf.MysqlAvatarRead["database"],
		conf.MysqlAvatarRead["charset"],
		conf.MysqlAvatarRead["parseTime"],
		conf.MysqlAvatarRead["loc"],
	)
	dbs["read"], err = gorm.Open("mysql", dsnRead)
	dsnWrite := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		conf.MysqlAvatarWrite["username"],
		conf.MysqlAvatarWrite["password"],
		conf.MysqlAvatarWrite["ip"],
		conf.MysqlAvatarWrite["port"],
		conf.MysqlAvatarWrite["database"],
		conf.MysqlAvatarWrite["charset"],
		conf.MysqlAvatarWrite["parseTime"],
		conf.MysqlAvatarWrite["loc"],
	)
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

	logFile := conf.MysqlAvatarLog[key]
	f, err := os.Create(logFile)
	if err != nil {
		fmt.Printf("get form err: %s", err.Error())
	}

	db.SetLogger(log.New(f, "\r\n", 0))
}
