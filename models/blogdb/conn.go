package blogdb

import (
	"gin-blog/interfaces"

	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
)

var DB *gorm.DB
var H interfaces.Helper

func Conn(h interfaces.Helper) (db *gorm.DB, err error) {
	H = h

	dbDsn := h.GetConfig("db_dsn")
	db, err = gorm.Open("sqlite3", dbDsn)
	if err != nil {
		return nil, err
	}

	// 连接成功
	DB = db
	DB.LogMode(true)
	SetGormLogger()
	DB.AutoMigrate(&SiteConfig{}, &Article{}, &Category{}, &Link{}) // 无表时会自动创建，有表时不会执行什么，除非实际表字段缺失
	LoadSiteConfig()                                                // 载入 site config
	return
}

// 日志输出到文件
func SetGormLogger() {
	logFile := H.GetConfig("blog_db_log_file").(string)
	f, err := os.Create(logFile)
	if err != nil {
		fmt.Printf("get form err: %s", err.Error())
	}
	DB.SetLogger(log.New(f, "\r\n", 0))
}
