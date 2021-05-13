package model

import (
	. "litrocket/common"
	"litrocket/utils/errmsg"
	"litrocket/utils/handlelog"
)

//Group 为数据库表Group的模型
type Group struct {
	GroupID   UserID `gorm:"type:int;not null"`
	UserRole  int    `gorm:"type:int;not null"` // 0 is manager
	UserID    UserID `gorm:"type:int;not null"`
	UserState int    `gorm:"type:int;not null"` // 0 is waiting to be added.
}

// Create A New Group.
func CreateGroup(group Group) {
	result := Db.Create(group)
	if result.Error != nil {
		handlelog.Handlelog("WARNING", "CreateGroup"+result.Error.Error())
	}
}

// Search Group's UserID And Number of Group's user.
func SearchGroupUser(groupusers []UserID, DestID UserID) int64 {
	var (
		i     int64
		num   int64
		users []Group
	)

	Db.Where("group_id = ?", DestID).Find(&users)
	num = int64(len(users))

	for i = 0; i < num; i++ {
		groupusers[i] = users[i].UserID
	}

	return num
}

// Join A Group
// todo Join A Group Need Group Manager's Agree.
func JoinGroup(userid, groupid UserID) int {
	group := Group{GroupID: groupid, UserRole: 1, UserID: userid, UserState: 0}
	if result := Db.Create(group); result.Error != nil {
		handlelog.Handlelog("WARNING", "JoinGroup + Db.Create"+result.Error.Error())
		return -1
	}

	return errmsg.OK_SUCCESS
}

// Get Some All Group Of Some User.
func GetAllGroup(userid UserID) []Group {
	var groups []Group

	Db.Where("user_id = ?", userid).Find(&groups)

	return groups
}

// Quit A Group
func QuitGroup(userid, groupid UserID) int {
	var group Group

	//* Quit User Is Manager.
	Db.Where("user_id = ? AND group_id = ?", userid, groupid).First(&group)
	if group.UserRole == 0 {
		result := Db.Where("group_id = ?", groupid).Delete(Group{})
		DelGroupInfo(groupid)
		if result.Error != nil {
			handlelog.Handlelog("WARNING", "QuitGroup"+result.Error.Error())
			return -1
		}
		return errmsg.OK_SUCCESS
	}

	//* Quit User Is Not Manager.
	result := Db.Where("user_id = ? AND group_id = ?", userid, groupid).Delete(&group)
	if result.Error != nil {
		handlelog.Handlelog("WARNING", "QuitGroup"+result.Error.Error())
		return -1
	}

	return errmsg.OK_SUCCESS
}
