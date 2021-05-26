package model

import (
	. "litrocket/common"
	"litrocket/utils/handlelog"
	"time"
)

//Message 为数据库表Message的模型
type Message struct {
	MsgDate    string `gorm:"type:datetime;not null"`
	Message    string `gorm:"type:longtext;not null"` // 文本消息或文件消息的JSON消息(直接把JSON存到数据库,发送的时候好发)
	SrcID      UserID `gorm:"type:int;not null"`
	DestID     UserID `gorm:"type:int;not null"`
	MessDest   int    `gorm:"type:int;not null"`
	MessFormat int    `gorm:"type:int;not null"` //消息类型, 0文本,1文件
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

// 将离线消息发送给用户
func SendOffLineMess(DestID UserID) []Message {
	var mes []Message
	Db.Where("dest_id = ?", DestID).Find(&mes)
	return mes
}
