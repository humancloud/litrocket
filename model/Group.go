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

func CreateGroup(group Group) {
	result := Db.Create(group)
	if result.Error != nil {
		handlelog.Handlelog("WARNING", "CreateGroup"+result.Error.Error())
	}
}

// SearchGroupUser
func SearchGroupUser(groupusers []UserID, DestID UserID) int64 {
	var (
		i     int64
		num   int64
		users []Group
	)

	result := Db.Where("group_id = ?", DestID).Find(&users)
	num = result.RowsAffected

	for i = 0; i < num; i++ {
		groupusers[i] = users[i].UserID
	}

	return num
}

// AddGroup
func AddGroup(userid, groupid UserID) int {
	//var groups []Group
	// 用户已经在群聊中
	//result := Db.Where("group_id = ? AND user_id = ?", groupid, userid).Find(&groups)
	//if result.RowsAffected > 0 {
	//	return errmsg.ERR_ALEALDY_GROUP
	//}

	group := Group{GroupID: groupid, UserRole: 1, UserID: userid, UserState: 0}
	if result := Db.Create(group); result.RowsAffected != 1 || result.Error != nil {
		handlelog.Handlelog("WARNING", "AddGroup+Db.Create"+result.Error.Error())
		return -1
	}

	return errmsg.OK_SUCCESS
}

// AddGroupok, i am the group's manager, i agree somebody's request.
func AddGroupOk(userid, groupid UserID) int {
	group := Group{UserState: 1}
	result := Db.Model(&group).Where("group_id = ? AND user_id = ?", groupid, userid).Update("user_state")
	if result.RowsAffected != 1 || result.Error != nil {
		handlelog.Handlelog("WARNING", "AddGroupOk+Update"+result.Error.Error())
		return -1
	}
	return errmsg.OK_SUCCESS
}

// GetAllGroup
func GetAllGroup(userid UserID) []Group {
	var groups []Group

	Db.Where("user_id = ?", userid).Find(&groups)

	return groups
}

// DelGroup()
func DelGroup(userid, groupid UserID) int {
	var group Group

	// user is manager.
	result := Db.Where("user_id = ? AND group_id = ?", userid, groupid).First(&group)
	if group.UserRole == 0 {
		result = Db.Where("group_id = ?", groupid).Delete(Group{})
		DelGroupInfo(groupid) //删除群信息
		if result.Error != nil {
			handlelog.Handlelog("WARNING", "DelGroup"+result.Error.Error())
			return -1
		}
		return errmsg.OK_SUCCESS
	}

	// user is not manager.
	result = Db.Where("user_id = ? AND group_id = ?", userid, groupid).Delete(&group)
	if result.RowsAffected != 1 || result.Error != nil {
		handlelog.Handlelog("WARNING", "DelGroup"+result.Error.Error())
		return errmsg.ERR_NOSUCHGROUP
	}

	return errmsg.OK_SUCCESS
}
