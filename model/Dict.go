package model

import (
	"litrocket/common"
	"litrocket/utils/errmsg"
	"litrocket/utils/handlelog"
)

type Dict struct {
	UserID  common.UserID `gorm:"type:int;not null"` // 0为本项目自带的单词,不为0为用户自己扩展的单词
	Chinese string        `gorm:"type:varchar(255);not null"`
	English string        `gorm:"type:varchar(255);not null"`
}

const MANAGER_UPLOAD_WORD = 0

func ChiLikeSearch(chinese string, id common.UserID) []Dict {
	var dicts []Dict

	Db.Where("chinese LIKE ? AND (user_id = ? OR user_id = ?)", "%"+chinese+"%", MANAGER_UPLOAD_WORD, id).Find(&dicts)

	return dicts
}

func EngSearch(english string, id common.UserID) Dict {
	var dict Dict

	Db.Where("english = ? AND (user_id = ? OR user_id = ?)", english, MANAGER_UPLOAD_WORD, id).First(&dict)

	return dict
}

func PushWord(chinese, english string, id common.UserID) int {
	var dicts []Dict
	dict := Dict{Chinese: chinese, English: english, UserID: id}

	// 已有单词就不要再录入了
	Db.Where("english = ? AND (user_id = ? OR user_id = ?)", english, MANAGER_UPLOAD_WORD, id).Find(&dicts)
	if len(dicts) > 0 {
		return errmsg.ERR_DICT_PUSH_REPEAT
	}

	result := Db.Create(&dict)
	if result.RowsAffected != 1 || result.Error != nil {
		handlelog.Handlelog("WARNING", "model:InsertUser Db.Create"+chinese+english+result.Error.Error())
		return -1
	}
	return errmsg.OK_SUCCESS
}
