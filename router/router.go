package router

import (
	apiv1 "litrocket/api/v1"
)

const (
	// user
	GetUser          = "get/userinfo"
	SearchUserByName = "sea/user"
	EditUserImage    = "edit/userimage"
	EditUserAge      = "edit/userage"
	EditUserSex      = "edit/usersex"
	EditUserTips     = "edit/usertips"
	// friend
	AddFriend     = "add/friend"
	AddFriendOk   = "add/friendok"
	GetAllFrend   = "get/friend"
	GetFriendInfo = "get/friendinfo"
	DeleteFriend  = "del/friend"
	RecommFrend   = "rec/friend"
	// group
	AddGroup          = "add/group"
	AddGroupOk        = "add/groupok"
	CreGroup          = "cre/group"
	GetAllGroup       = "get/group"
	GetGroupInfo      = "get/groupinfo"
	DeleteGroup       = "del/group"
	UploadGroupImage  = "upload/groupimg"
	SearchGroupByName = "sea/group"
	EditGroupTips     = "edit/grouptips"
	EditGroupImage    = "edit/groupimage"
	// Chat Send message
	SendMess = "send/mess"
	// My News
	MyNews  = "upload/mynews"
	GetNews = "get/news"
	// File
	PerSendFile   = "personsend/file"
	PerRecvFile   = "personrecv/file"
	UpGroupFile   = "upgroup/file"
	DownGroupFile = "downgroup/file"
	ViewGroupFile = "viewgroup/file"
	// Video
	JoinVideo = "join/screen"
	QuitVideo = "quit/video"
	SendVideo = "send/screen"
	EndVideo  = "end/video"
	SeaVideo  = "sea/video"
)

var routerv1 = make(map[string]func(json []byte))

func InitRouter() {
	// user.
	routerv1[GetUser] = apiv1.GetUserInfo
	routerv1[SearchUserByName] = apiv1.SearchUserByName
	routerv1[EditUserImage] = apiv1.EditUserImage
	routerv1[EditUserAge] = apiv1.EditUserAge
	routerv1[EditUserSex] = apiv1.EditUserSex
	routerv1[EditUserTips] = apiv1.EditUserTips
	//friend.
	routerv1[AddFriend] = apiv1.AddFriend
	routerv1[AddFriendOk] = apiv1.Agree
	routerv1[GetAllFrend] = apiv1.GetAllFriend
	routerv1[GetFriendInfo] = apiv1.GetFriendInfo
	routerv1[DeleteFriend] = apiv1.DelFriend
	routerv1[RecommFrend] = apiv1.FriendRecommand
	//group.
	routerv1[AddGroup] = apiv1.AddGroup
	routerv1[AddGroupOk] = apiv1.AddGroupOk
	routerv1[CreGroup] = apiv1.CreateGroup
	routerv1[GetAllGroup] = apiv1.GetAllGroup
	routerv1[GetGroupInfo] = apiv1.GetGroupInfo
	routerv1[DeleteGroup] = apiv1.DelGroup
	routerv1[UploadGroupImage] = apiv1.UploadGroupImage
	routerv1[SearchGroupByName] = apiv1.SearchGroupByName
	routerv1[EditGroupTips] = apiv1.EditGroupTips
	routerv1[EditGroupImage] = apiv1.EditGroupImage
	//Chat send mess.
	routerv1[SendMess] = apiv1.ChatServe
	//file
	routerv1[PerSendFile] = apiv1.SendPersonalFile
	routerv1[PerRecvFile] = apiv1.RecvPersonalFile
	routerv1[UpGroupFile] = apiv1.UploadGroupFile
	routerv1[DownGroupFile] = apiv1.DownloadGroupFile
	routerv1[ViewGroupFile] = apiv1.ViewGroupFiles
	//news
	//video
	routerv1[JoinVideo] = apiv1.JoinScreen
	routerv1[QuitVideo] = apiv1.QuitScreen
	routerv1[SendVideo] = apiv1.SendScreen
	routerv1[EndVideo] = apiv1.EndScreen
	routerv1[SeaVideo] = apiv1.SeaVideo
}

func Run(url string, json []byte) {
	if fuc, ok := routerv1[url]; ok {
		go fuc(json)
	}
}
