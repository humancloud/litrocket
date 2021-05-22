package model

import (
	. "litrocket/common"
	"time"

	"github.com/jinzhu/gorm"
)

//Article 为数据库表Article的模型,存储用户动态
type Article struct {
	gorm.Model
	UserID  UserID `gorm:"type:int;not null"`
	Time    string `gorm:"type:datetime;not null"`
	Content string `gorm:"type:varchar(255);not null"`
}

func CreateArticle(location string, userid UserID) (bool, uint) {
	art := Article{UserID: userid, Content: location, Time: time.Now().Format("2006-01-02 15:04:05")}
	result := Db.Create(&art)
	return (result.Error == nil), art.ID
}

func SearchArticle(ArticleId, UserId UserID) bool {
	var art []Article
	Db.Where("user_id = ? AND id = ?", UserId, ArticleId).Find(&art)
	return (len(art) > 0)
}

func SearchArticleById(ArticleId, UserId UserID) []Article {
	var art []Article
	Db.Where("user_id = ? AND id = ?", UserId, ArticleId).First(&art)
	return art
}

func GetAllArticle(UserId UserID) []Article {
	var art []Article
	Db.Where("user_id = ?", UserId).Find(&art)
	return art
}

func DelArticle(ArticleId, UserId UserID) (int, string) {
	var art []Article
	Db.Where("user_id = ? AND id = ?", UserId, ArticleId).Find(&art)
	Db.Where("user_id = ? AND id = ?", UserId, ArticleId).Delete(Article{})
	return int(art[0].ID), art[0].Content
}
