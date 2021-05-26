package router

import (
	"fmt"
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
	DelGroupFile  = "delgroup/file"
	// Video
	JoinVideo = "join/screen"
	QuitVideo = "quit/video"
	SendVideo = "send/screen"
	EndVideo  = "end/video"
	SeaVideo  = "sea/video"
	// PersonalDict
	PersonalChiDict = "sea/chidict"
	PersonalEnglish = "sea/engdict"
	PushWord        = "push/word"
	// Offline mess
	OffLineMess = "offline/mess"
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
	//routerv1[AddFriendOk] = apiv1.Agree
	routerv1[GetAllFrend] = apiv1.GetAllFriend
	routerv1[GetFriendInfo] = apiv1.GetFriendInfo
	routerv1[DeleteFriend] = apiv1.DelFriend
	routerv1[RecommFrend] = apiv1.FriendRecommand
	//group.
	routerv1[AddGroup] = apiv1.AddGroup
	//routerv1[AddGroupOk] = apiv1.AddGroupOk
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
	routerv1[DelGroupFile] = apiv1.DeleteGroupFile
	//news
	//video
	routerv1[JoinVideo] = apiv1.JoinScreen
	routerv1[QuitVideo] = apiv1.QuitScreen
	routerv1[SendVideo] = apiv1.SendScreen
	routerv1[EndVideo] = apiv1.EndScreen
	routerv1[SeaVideo] = apiv1.SeaVideo
	// Dict
	routerv1[PersonalChiDict] = apiv1.PersonalChiDict
	routerv1[PersonalEnglish] = apiv1.PersonalEngDict
	routerv1[PushWord] = apiv1.PushWord
	// OfflineMess
	routerv1[OffLineMess] = apiv1.SendOffLineMess
}

func Run(url string, json []byte) {
	fmt.Println(string(json))
	if fuc, ok := routerv1[url]; ok {
		fuc(json)
	}

	//! API不能直接创建新协程处理,若用户连续多个请求,服务端又同时向客户端返回数据,客户端一个Resp连接接收必然出错.
	//! File 单线程的话,时间很久,绝不行.
	//! Video 目前是调用另外的一个应用处理媒体数据转发,所以可以不用新协程.
	//! Chat 多协程的话,向一个用户连续发两条消息,两个协程同时给对方发送,一定出错.
	//* 因此: Chat还是单协程,因为把chat另起一个应用的话时间也并不会更快.
	//* Video 如果不另起一个应用,那就在videoAPI里面创建协程处理.
	//* File  搭建一个FTP应用.
	//* 其他API处理速度较快,都是单协程.
	//* 像File,Video 后续要加的东西较多,且处理时间可能会很长,因此不应算作API里面,应当是属于额外模块.   Chat要加图片消息,语音消息,其实也没什么好加的,Chat还可以继续在API里面,最好是独立为模块
}
