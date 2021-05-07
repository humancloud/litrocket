package model

import "litrocket/common"

//Article 为数据库表Article的模型,存储用户动态
type Article struct {
	UserID  common.UserID `gorm:"type:int;not null"`
	Time    string        `gorm:"type:datetime;not null"`
	Content string        `gorm:"type:longtext;not null"`
	Image   string        `gorm:"type:varchar(100);not null"`
}
