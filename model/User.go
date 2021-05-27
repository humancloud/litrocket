package model

import (
	. "litrocket/common"
	"litrocket/utils/errmsg"
	"litrocket/utils/handlelog"

	"github.com/jinzhu/gorm"
)

const (
	USERHEADIMG = "tempfile/headimg/personal/"
)

//User 为数据库表User的模型
type User struct {
	gorm.Model
	UserName  string `gorm:"type:varchar(255);not null"`
	PassWord  string `gorm:"type:varchar(255);not null"`
	UserMail  string `gorm:"type:varchar(255);not null"`
	UserSex   int    `gorm:"type:int"`
	UserAge   string `gorm:"type:varchar(10)"`
	UserImage string `gorm:"type:longtext"`
	UserTips  string `gorm:"type:varchar(4096)"`
}

// Search User By Name.
// Return Passwd,ID,exist.
func SearchUser(Name string) (User, bool) {
	var users []User
	var temp User
	Db.Where("user_name = ?", Name).Find(&users) // 查询时UserName变成了user_name, 因为数据库字段是被gorm修改为这样的

	if len(users) != 1 {
		return temp, false
	}
	return users[0], true
}

// Search User By ID
// Return User{}, exist.
func SearchByID(ID UserID) (User, bool) {
	var (
		users []User
		temp  User
	)
	Db.Where("id = ?", ID).Find(&users)
	if len(users) != 1 {
		return temp, false
	}

	return users[0], true
}

// Mail is exist ?
func SearchMail(Mail string) bool {
	var users []User

	Db.Where("user_mail = ?", Mail).Find(&users)
	if len(users) != 1 {
		return false
	}
	return true
}

// InsertUser
func InsertUser(Name, Passwd, Mail string) bool {
	user := User{UserName: Name, PassWord: Passwd, UserMail: Mail}

	result := Db.Create(&user)
	if result.RowsAffected != 1 || result.Error != nil {
		handlelog.Handlelog("WARNING", "model:InsertUser Db.Create"+Name+Passwd+Mail+result.Error.Error())
		return false
	}

	return true
}

// Search Users where username like "Name"
func SearchByName(Name string) []User {
	var users []User

	Db.Where("user_name LIKE ?", "%"+Name+"%").Find(&users)

	return users
}

func UpDatePasswd(passwd, name string) int {
	result := Db.Model(&User{}).Where("user_name = ?", name).Update("pass_word", passwd)
	if result.Error != nil {
		handlelog.Handlelog("WARNING", "UpDatePasswd"+result.Error.Error())
		return errmsg.ERR_UNKNOWN
	}
	return errmsg.OK_SUCCESS
}

func UploadUserImage(img string, ID UserID) int {
	result := Db.Model(&User{}).Where("id = ?", ID).Update("user_image", img)
	if result.Error != nil {
		handlelog.Handlelog("WARNING", "UploadUserImage"+result.Error.Error())
		return -1
	}

	return errmsg.OK_SUCCESS
}

func EditUserAge(Id UserID, age int) int {
	result := Db.Model(&User{}).Where("id = ?", Id).Update("user_age", age)
	if result.Error != nil {
		handlelog.Handlelog("WARNING", "EditUserAge"+result.Error.Error())
		return -1
	}
	return errmsg.OK_SUCCESS
}

func EditUserSex(Id UserID, sex int) int {
	result := Db.Model(&User{}).Where("id = ?", Id).Update("user_sex", sex)
	if result.Error != nil {
		handlelog.Handlelog("WARNING", "EditUserSex"+result.Error.Error())
		return -1
	}
	return errmsg.OK_SUCCESS
}

func EditUserTips(Id UserID, tips string) int {
	result := Db.Model(&User{}).Where("id = ?", Id).Update("user_tips", tips)
	if result.Error != nil {
		handlelog.Handlelog("WARNING", "EditUserTips"+result.Error.Error())
		return -1
	}
	return errmsg.OK_SUCCESS
}
