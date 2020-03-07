package blogdb

import "github.com/jinzhu/gorm"

type SiteConfig struct {
	gorm.Model

	ConfigName  string `gorm:"type:varchar(100);unique_index;not null"`
	ConfigValue string `gorm:"type:varchar(255);"`
}

var siteConfigMap = make(map[string]interface{})

// 加载数据库的 web 配置
// 每次后台修改过配置时，都应调用此方法
func LoadSiteConfig() {
	var siteConfig []*SiteConfig
	DB.Select("config_name,config_value").Find(&siteConfig)

	for _, item := range siteConfig {
		siteConfigMap[item.ConfigName] = item.ConfigValue
	}
}

// 获取配置
func GetSiteConfig(key string) interface{} {
	return siteConfigMap[key]
}
