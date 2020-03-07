package conf

var (
	Port        = ":8000"
	Domain      = ""
	Debug       = true
	AppPath     = ""
	RuntimePath = "runtime/"

	BlogDB    = "blog.db?_loc=Asia/Shanghai"
	BlogDBLog = RuntimePath + "logs/blogdb.log"

	LogCommonFileForGin = RuntimePath + "logs/app_common_gin.log" // gin 默认输出的日志
	LogCommonFile       = RuntimePath + "logs/app_common.log"     // 代码级 error，意料内的 error
	LogErrorFile        = RuntimePath + "logs/app_error.log"      // recover 的 error，意料之外，由 gin 捕获

	MysqlAvatarRead = map[string]string{
		"username":  "root",
		"password":  "root",
		"ip":        "127.0.0.1",
		"database":  "spider_coo_avatar",
		"charset":   "utf8",
		"parseTime": "True",
		"loc":       "Local",
		"port":      "3306",
	}
	MysqlAvatarWrite = map[string]string{
		"username":  "root",
		"password":  "root",
		"ip":        "127.0.0.1",
		"database":  "spider_coo_avatar",
		"charset":   "utf8",
		"parseTime": "True",
		"loc":       "Local",
		"port":      "3306",
	}
	MysqlAvatarLog = map[string]string{
		"read":  RuntimePath + "logs/avatar_read.log",
		"write": RuntimePath + "logs/avatar_write.log",
	}

	Beanstalkd = "127.0.0.1:11300"
)
