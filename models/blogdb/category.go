package blogdb

import "github.com/jinzhu/gorm"

type Category struct {
	gorm.Model

	Name         string     `gorm:"type:vachar(50)"`
	Articles     []*Article `gorm:"-"`
	ArticleTotal uint       `gorm:"type:int"`
	Visits       uint
	ArticleCount uint
}

func (c *Category) GetList(page, num uint64, where string, args ...interface{}) (categories []*Category) {
	if where == "" {
		where = "id > 0"
	}
	DB.Order("visits desc").Where(where, args...).Offset((page - 1) * num).Limit(num).Find(&categories)

	return
}

func (c *Category) Count(where string, args ...interface{}) (count uint64) {
	if where == "" {
		where = "id > 0"
	}
	DB.Model(&Category{}).Where(where, args...).Select("id").Count(&count)
	return
}

func (c *Category) PageCount(num uint64, where string, args ...interface{}) uint64 {
	count := c.Count(where, args...)
	p := count % num
	if p*num == count {
		return p
	} else {
		return p + 1
	}
}

func (c *Category) First(where string, args ...interface{}) *Category {
	if where == "" {
		DB.First(c)
	} else {
		DB.Where(where, args...).First(c)
	}

	return c
}
