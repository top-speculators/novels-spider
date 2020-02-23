package blogdb

import "github.com/jinzhu/gorm"

type Article struct {
	gorm.Model

	Title      string `gorm:"type:varchar(200)"`
	Content    string `gorm:"type:text"`
	Visits     uint   `gorm:"type:int"`
	CategoryId uint   `gorm:"type:int"`
	Category   Category
}

func (a *Article) GetList(page, num uint64, where string, args ...interface{}) (articles []*Article) {
	if where == "" {
		where = "id > 0"
	}
	selectQuery := "id,created_at,category_id,title,visits,SUBSTR(content,0,INSTR(content,'<!--more-->')) as content"
	order := "created_at desc"
	offset := (page - 1) * num
	DB.Select(selectQuery).Where(where, args...).Order(order).Offset(offset).Limit(num).Find(&articles)
	return
}

func (a *Article) Count(where string, args ...interface{}) (count uint64) {
	if where == "" {
		where = "id > 0"
	}
	DB.Model(&Article{}).Where(where, args...).Select("id").Count(&count)
	return
}

func (a *Article) PageCount(num uint64, where string, args ...interface{}) uint64 {
	count := a.Count(where, args...)
	p := count % num
	if p*num == count {
		return p
	} else {
		return p + 1
	}
}

func (a *Article) First(where string, args ...interface{}) *Article {
	if where == "" {
		DB.First(a)
	} else {
		DB.Where(where, args...).First(a)
	}

	return a
}
