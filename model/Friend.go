package model

import (
	. "litrocket/common"
	"litrocket/utils/errmsg"
	"litrocket/utils/handlelog"
)

//Friend 是数据库表Friend的模型
type Friend struct {
	UserID      UserID `gorm:"type:int;not null"`
	FriendID    UserID `gorm:"type:int;not null"`
	FriendState int    `gorm:"type:int;not null"` // 0 暂时未同意， 1 已经成为好友
	FriendNotes string `gorm:"type:varchar(255)"` //备注
}

// AddFriend
func AddFriend(MyID, FriendID UserID) int {
	var users []User
	var friends []Friend

	// Search User exist ?
	result := Db.Where("id = ?", FriendID).Find(&users)
	if result.RowsAffected != 1 {
		return errmsg.ERR_FRIEND_NO_EXIST
	}

	// Search we already are friends.
	result = Db.Where("user_id = ? AND friend_id = ?", MyID, FriendID).Find(&friends)
	if result.RowsAffected > 0 {
		return errmsg.ERR_ALEALDY_FRIEND
	}

	// Create
	friend := Friend{MyID, FriendID, 0, ""}
	result = Db.Create(&friend)
	if result.RowsAffected != 1 || result.Error != nil {
		handlelog.Handlelog("WARNING", "AddFriend() :"+result.Error.Error())
		return errmsg.ERR_UNKNOWN
	}

	friend = Friend{FriendID, MyID, 0, ""}
	result = Db.Create(&friend)
	if result.RowsAffected != 1 || result.Error != nil {
		handlelog.Handlelog("WARNING", "AddFriend() :"+result.Error.Error())
		return errmsg.ERR_UNKNOWN
	}

	return errmsg.OK_SUCCESS
}

// I Agree friend request.
func AddFriendOk(MyID, FriendID UserID) int {
	friend := Friend{FriendState: 1}

	result := Db.Model(&friend).Where("user_id = ? AND friend_id = ?", MyID, FriendID).Update("friend_state")
	if result.RowsAffected != 1 || result.Error != nil {
		handlelog.Handlelog("WARNING", "AddFriendOk() :"+result.Error.Error())
		return errmsg.ERR_UNKNOWN
	}

	result = Db.Model(&friend).Where("user_id = ? AND friend_id = ?", FriendID, MyID).Update("friend_state")
	if result.RowsAffected != 1 || result.Error != nil {
		handlelog.Handlelog("WARNING", "AddFriendOk() :"+result.Error.Error())
		return errmsg.ERR_UNKNOWN
	}

	return errmsg.OK_SUCCESS
}

// Get All My Friend.
func GetAllFriend(MyID UserID) []Friend {
	var friends []Friend

	Db.Where("user_id = ?", MyID).Find(&friends)

	return friends
}

// Delete Friend
func DelFriend(MyID, FriendID UserID) int {
	// Delete
	result := Db.Where("friend_id = ? AND user_id = ?", MyID, FriendID).Delete(Friend{})
	result = Db.Where("friend_id = ? AND user_id = ?", FriendID, MyID).Delete(Friend{})
	if result.RowsAffected != 1 || result.Error != nil {
		handlelog.Handlelog("WARNING", "DelFriend() :"+result.Error.Error())
		return errmsg.ERR_UNKNOWN
	}

	return errmsg.OK_SUCCESS
}

// EditFriendNotes
func EditFriendNotes(MyID, FriendID UserID, Notes string) int {
	friend := Friend{FriendNotes: Notes}
	result := Db.Model(&friend).Where("user_id = ? AND friend_id = ?", MyID, FriendID).Update("friend_notes")
	if result.RowsAffected != 1 || result.Error != nil {
		handlelog.Handlelog("WARNING", "EditFriendNotes :"+result.Error.Error())
		return errmsg.ERR_UNKNOWN
	}

	return errmsg.OK_SUCCESS
}
