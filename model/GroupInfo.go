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

//GroupInfo 为数据库表GroupInfo的模型
type GroupInfo struct {
	//model中已经包含，ID,createtime,updatetime,deletetime字段
	gorm.Model
	GroupName    string `gorm:"type:varchar(20);not null"`
	GroupImage   string `gorm:"type:varchar(255)"`
	GroupTips    string `gorm:"type:varchar(4096)"`
	GroupRootID  UserID `gorm:"type:int;not null"` // 群聊创建者
	GroupUserNum int    `gorm:"type:int;not null"` // (人数达到200为满)
}

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

func DelGroupInfo(groupid UserID) {
	result := Db.Where("id = ?", groupid).Delete(GroupInfo{})
	if result.Error != nil {
		handlelog.Handlelog("WARNING", "DelGroupInfo"+result.Error.Error())
	}
}

func SearchGroupExist(name string) bool {
	var g []GroupInfo
	Db.Where("group_name = ?", name).Find(&g)
	if len(g) > 0 {
		return true
	}

	return false
}

func SearchGroupIdByName(name string) UserID {
	var group GroupInfo
	Db.Where("group_name = ?", name).First(&group)
	return UserID(group.ID)
}

// 不向客户端发送群聊头像，那样JSON太大
// 发送结构体内的除GroupImage以外的所有东西
func SearchGroupByName(name string) []GroupInfo {
	var groups []GroupInfo
	result := Db.Where("group_name LIKE ?", "%"+name+"%").Find(&groups)
	if result.Error != nil {
		handlelog.Handlelog("WARNING", "SearchGroupByName"+result.Error.Error())
	}
	return groups
}

func GetGroupInfo(GroupID UserID) GroupInfo {
	var group GroupInfo
	Db.Where("id = ?", GroupID).First(&group)
	return group
}

func EditGroupTips(ID UserID, Tips *string) int {
	result := Db.Model(&GroupInfo{}).Where("id = ?", ID).Update("group_tips", *Tips)
	if result.Error != nil {
		handlelog.Handlelog("WARNING", "EditGroupTips"+result.Error.Error())
		return -1
	}
	return errmsg.OK_SUCCESS
}

func GetGroupImg(ID UserID) string {
	var group GroupInfo
	Db.Where("id = ?", ID).First(&group)
	return group.GroupImage
}

func UploadGroupImage(img string, ID UserID) {
	Db.Model(&GroupInfo{}).Where("id = ?", ID).Update("group_image", img)
}
