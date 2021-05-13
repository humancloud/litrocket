package model

import (
	"litrocket/common"
	"litrocket/utils/errmsg"
	"litrocket/utils/handlelog"
)

const (
	PERSONFILE = 0
	GROUPFILE  = 1

	NOSEND  = 0
	STROAGE = 1
)

type File struct {
	Filetime     string        `gorm:"type:datetime;not null"`
	FileSize     int64         `gorm:"int;not null"`
	Filename     string        `gorm:"type:varchar(255);not null"`
	SrcName      string        `gorm:"type:varchar(255);not null"`
	FileLocation string        `gorm:"type:varchar(255);not null"`
	SrcID        common.UserID `gorm:"int;not null"`
	DestID       common.UserID `gorm:"int;not null"`
	FileType     int           `gorm:"int;not null"` // 个人文件0,群聊文件1
	FileState    int           `gorm:"int;not null"` // 未发送0(属于离线文件)，已发送1(属于群文件或个人历史发送文件)
}

func UploadGroupFile(file *File) {
	result := Db.Create(file)
	if result.Error != nil {
		handlelog.Handlelog("WARNING", "UploadGroupFile"+result.Error.Error())
	}
}

func ViewGroupFiles(groupid common.UserID) []File {
	var files []File

	Db.Where("dest_id = ? AND file_type = ?", groupid, GROUPFILE).Find(&files)

	return files
}

func DeleteGroupFile(filename string, srcid common.UserID, ID common.UserID) (int, string) {
	var file File
	result := Db.Where("filename = ? AND src_id = ? AND dest_id = ? AND file_type = ?", filename, srcid, ID, GROUPFILE).First(&file).Delete(File{})
	if result.Error != nil {
		handlelog.Handlelog("WARNING", "DeleteGroupFile"+result.Error.Error())
		return -1, ""
	}

	if result.RowsAffected != 1 {
		return -1, ""
	}

	return errmsg.OK_SUCCESS, file.FileLocation
}

func GetGroupFileLocation(filename string, ID common.UserID) (string, int64) {
	var file File
	Db.Where("filename = ? AND dest_id = ? AND file_type = ?", filename, ID, GROUPFILE).First(&file)
	return file.FileLocation, file.FileSize
}

func PersonFile(file *File) {
	result := Db.Create(file)
	if result.Error != nil {
		handlelog.Handlelog("WARNING", "PersonFile"+result.Error.Error())
	}
}

func GetFileLocation(filename string, SrcID, DestID common.UserID) (string, int64) {
	var file File
	Db.Where("filename = ? AND src_id = ? AND dest_id = ? AND file_type = ?", filename, SrcID, DestID, PERSONFILE).First(&file)
	return file.FileLocation, file.FileSize
}
