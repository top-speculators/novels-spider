package blogdb

import "github.com/jinzhu/gorm"

type Link struct {
	gorm.Model

	Name string `gorm:"type:varchar(20)"`
	Href string `gorm:"type:varchar(20)"`
}
