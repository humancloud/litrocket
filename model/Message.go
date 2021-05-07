package model

import (
	. "litrocket/common"
	"litrocket/utils/handlelog"
	"time"
)

//Message 为数据库表Message的模型
type Message struct {
	MsgDate    string `gorm:"type:datetime;not null"`
	Message    string `gorm:"type:longtext;not null"`
	SrcID      UserID `gorm:"type:int;not null"`
	DestID     UserID `gorm:"type:int;not null"`
	MessDest   int    `gorm:"type:int;not null"`
	MessFormat int    `gorm:"type:int;not null"`
}

// SaveMess 将消息保存到离线消息表
func SaveMess(mess *string, SrcID, DestID UserID, MessDest, MessFormat int) {
	Mess := Message{
		MsgDate:    time.Now().Format("2006-01-02 15:04:05"),
		Message:    (*mess),
		SrcID:      SrcID,
		DestID:     DestID,
		MessDest:   MessDest,
		MessFormat: MessFormat,
	}

	result := Db.Create(&Mess)
	if result.RowsAffected != 1 || result.Error != nil {
		handlelog.Handlelog("WARNING", "model:SaveMess Db.Create"+result.Error.Error())
	}
}
