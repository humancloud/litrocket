package model

import (
	. "litrocket/common"
	"litrocket/utils/errmsg"
	"litrocket/utils/handlelog"

	"github.com/jinzhu/gorm"
)

const (
	GROUPIMGDIR = "tempfile/headimg/group/"
)

type GroupInfo struct {
	gorm.Model          // model中已经包含，ID,createtime,updatetime,deletetime字段
	GroupName    string `gorm:"type:varchar(20);not null"`
	GroupImage   string `gorm:"type:varchar(255)"`
	GroupTips    string `gorm:"type:varchar(4096)"`
	GroupRootID  UserID `gorm:"type:int;not null"` // 群聊创建者
	GroupUserNum int    `gorm:"type:int;not null"` // (人数达到200为满)
}

// Create A New Group And It's Info.
func CreateGroupInfo(group GroupInfo) UserID {
	var g GroupInfo
	result := Db.Create(&group)
	if result.Error != nil {
		handlelog.Handlelog("WARNING", "CreateGroupInfo,Create"+result.Error.Error())
	}
	result = Db.Where("group_name = ?", group.GroupName).First(&g)
	if result.Error != nil {
		handlelog.Handlelog("WARNING", "CreateGroupInfo,Where"+result.Error.Error())
	}
	return UserID(g.ID)
}

// Delete A Group And It's Info.
func DelGroupInfo(groupid UserID) {
	result := Db.Where("id = ?", groupid).Delete(GroupInfo{})
	if result.Error != nil {
		handlelog.Handlelog("WARNING", "DelGroupInfo"+result.Error.Error())
	}
}

// Group Is Exist ?
func SearchGroupExist(name string) bool {
	var g []GroupInfo

	Db.Where("group_name = ?", name).Find(&g)
	return len(g) > 0
}

// Get Group's ID By Search Group Name.
func SearchGroupIdByName(name string) UserID {
	var group GroupInfo
	Db.Where("group_name = ?", name).First(&group)
	return UserID(group.ID)
}

// Search GroupName By LIKE Name.
func SearchGroupByName(name string) []GroupInfo {
	var groups []GroupInfo
	result := Db.Where("group_name LIKE ?", "%"+name+"%").Find(&groups)
	if result.Error != nil {
		handlelog.Handlelog("WARNING", "SearchGroupByName"+result.Error.Error())
	}
	return groups
}

// Get Some Group's Info.
func GetGroupInfo(GroupID UserID) GroupInfo {
	var group GroupInfo
	Db.Where("id = ?", GroupID).First(&group)
	return group
}

// Edit Group's tips.
func EditGroupTips(ID UserID, Tips *string) int {
	result := Db.Model(&GroupInfo{}).Where("id = ?", ID).Update("group_tips", *Tips)
	if result.Error != nil {
		handlelog.Handlelog("WARNING", "EditGroupTips"+result.Error.Error())
		return -1
	}
	return errmsg.OK_SUCCESS
}

// Get Group's HeadImg.
func GetGroupImg(ID UserID) string {
	var group GroupInfo
	Db.Where("id = ?", ID).First(&group)
	return group.GroupImage
}

// Upload Group's HeadImg.
func UploadGroupImage(img string, ID UserID) {
	Db.Model(&GroupInfo{}).Where("id = ?", ID).Update("group_image", img)
}
