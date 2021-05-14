package apiv1

import (
	"io"
	. "litrocket/common"
	"litrocket/model"
	"litrocket/utils/dataencry"
	"litrocket/utils/errmsg"
	"litrocket/utils/handlelog"
	"os"
	"strconv"
)

type Group struct {
	Url       string
	MyID      UserID
	GroupID   UserID
	GroupName string
	Mess      string
}

type GroupResult struct {
	Code    int
	Url     string
	GroupID UserID
}

// Create A New Group.
func CreateGroup(json []byte) {
	var (
		G     model.GroupInfo
		id    UserID
		g     model.Group
		group struct {
			Url       string
			RootID    UserID
			GroupName string
		}

		result struct {
			Url   string
			Code  int
			Group model.GroupInfo
		}
	)

	if err := dataencry.Unmarshal(json, &group); err != nil {
		return
	}

	result.Url = group.Url

	// Group Is Already Exist ?
	if ok := model.SearchGroupExist(group.GroupName); ok {
		result.Code = -1
		goto WRITE
	}

	//* Create A New Group.
	G = model.GroupInfo{GroupName: group.GroupName, GroupRootID: group.RootID, GroupUserNum: 1}
	id = model.CreateGroupInfo(G)
	g = model.Group{GroupID: id, UserRole: 0, UserID: group.RootID, UserState: 0}
	model.CreateGroup(g)

	result.Code = errmsg.OK_SUCCESS
	result.Group.ID = uint(id)
	result.Group.GroupName = group.GroupName

WRITE:
	r, _ := dataencry.Marshal(result)
	if conns, ok := AllUsers.Load(group.RootID); ok {
		conn := conns.(Conns)
		b := append(r, []byte("\r\n--\r\n")...)
		conn.ResponseConn.Write(b)
	}
}

func AddGroup(json []byte) {
	var (
		group  Group
		err    error
		result struct {
			Url   string
			Code  int
			Group model.GroupInfo
		}
	)

	if err = dataencry.Unmarshal(json, &group); err != nil {
		return
	}

	// 查询群聊id
	result.Url = "add/groupresult"
	result.Group, result.Code = model.SearchGroup(group.GroupName)
	if result.Code == errmsg.OK_SUCCESS {
		// 直接添加到数据库
		result.Code = model.JoinGroup(group.MyID, UserID(result.Group.ID))
	}

	b, _ := dataencry.Marshal(result)
	r := append(b, []byte("\r\n--\r\n")...)
	if conns, ok := AllUsers.Load(group.MyID); ok {
		conn := conns.(Conns)
		conn.ResponseConn.Write(r)
	}
}

func GetAllGroup(json []byte) {
	var (
		group  Group
		err    error
		result struct {
			Url    string
			Groups []model.Group
		}
	)

	if err = dataencry.Unmarshal(json, &group); err != nil {
		return
	}

	result.Groups = model.GetAllGroup(group.MyID)
	result.Url = group.Url

	r, _ := dataencry.Marshal(result)
	r = append(r, []byte("\r\n--\r\n")...)

	if conns, ok := AllUsers.Load(group.MyID); ok {
		conn := conns.(Conns)
		conn.ResponseConn.Write(r)
	}
}

func GetGroupInfo(json []byte) {
	var (
		group struct {
			Url     string
			MyID    UserID
			GroupID UserID
		}
	)

	var (
		result struct {
			Url  string
			Info model.GroupInfo
		}
	)

	if err := dataencry.Unmarshal(json, &group); err != nil {
		return
	}

	result.Info = model.GetGroupInfo(group.GroupID)
	result.Url = group.Url
	r, _ := dataencry.Marshal(result)
	r = append(r, []byte("\r\n--\r\n")...)

	if conns, ok := AllUsers.Load(group.MyID); ok {
		conn := conns.(Conns)
		conn.ResponseConn.Write(r)
	}
}

func DelGroup(json []byte) {
	var (
		i          int64
		group      Group
		err        error
		groupusers = make([]UserID, 200)
	)

	if err = dataencry.Unmarshal(json, &group); err != nil {
		return
	}

	// 查询群聊的人数,先查询后删除群聊
	num := model.SearchGroupUser(groupusers, group.GroupID)

	model.QuitGroup(group.MyID, group.GroupID)

	// 向在线的群成员发送xx退出群聊消息,或是群聊解散消息
	for i = 0; i < num; i++ {
		if group.MyID == groupusers[i] {
			continue
		}

		if conns, ok := AllUsers.Load(groupusers[i]); ok {
			conn := conns.(Conns)
			b := append(json, []byte("\r\n--\r\n")...)
			conn.ResponseConn.Write(b)
		}
	}
}

func UploadGroupImage(json []byte) {
	var (
		f      *os.File
		group  Group
		err    error
		result struct {
			Url  string
			Code int
		}
	)

	if err = dataencry.Unmarshal(json, &group); err != nil {
		return
	}

	result.Url = group.Url
	img := model.GROUPIMGDIR + strconv.Itoa(int(group.GroupID))
	if f, err = os.OpenFile(img, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
		handlelog.Handlelog("WARNING", "UploadGroupImage"+err.Error())
		result.Code = -1
		goto WRITE
	}

	result.Code = errmsg.OK_SUCCESS

	io.WriteString(f, group.Mess)

	// Save to DB.
	model.UploadGroupImage(img, group.GroupID)

WRITE:
	if conns, ok := AllUsers.Load(group.MyID); ok {
		conn := conns.(Conns)
		r, _ := dataencry.Marshal(result)
		conn.ResponseConn.Write(r)
	}
}

func EditGroupImage(json []byte) {
	var (
		f      *os.File
		group  Group
		err    error
		result struct {
			Url  string
			Code int
		}
	)

	if err = dataencry.Unmarshal(json, &group); err != nil {
		return
	}

	result.Url = group.Url

	location := model.GetGroupImg(group.GroupID)
	if f, err = os.Open(location); err != nil {
		result.Code = -1
		goto WRITE
	}

	result.Code = errmsg.OK_SUCCESS
	io.WriteString(f, group.Mess)

WRITE:
	if conns, ok := AllUsers.Load(group.MyID); ok {
		conn := conns.(Conns)
		r, _ := dataencry.Marshal(result)
		conn.ResponseConn.Write(r)
	}
}

func EditGroupTips(json []byte) {
	var (
		group  Group
		err    error
		result struct {
			Url  string
			Code int
		}
	)

	if err = dataencry.Unmarshal(json, &group); err != nil {
		return
	}

	result.Code = model.EditGroupTips(group.MyID, &group.Mess)
	result.Url = group.Url
	r, _ := dataencry.Marshal(result)

	if conns, ok := AllUsers.Load(group.MyID); ok {
		conn := conns.(Conns)
		conn.ResponseConn.Write(r)
	}
}

func SearchGroupByName(json []byte) {
	var (
		group  Group
		err    error
		result struct {
			Url    string
			Groups []string
		}
	)

	if err = dataencry.Unmarshal(json, &group); err != nil {
		return
	}

	Groups := model.SearchGroupByName(group.GroupName)
	for i := 0; i < len(Groups); i++ {
		result.Groups = append(result.Groups, Groups[i].GroupName)
	}
	result.Url = group.Url
	r, _ := dataencry.Marshal(result)

	if conns, ok := AllUsers.Load(group.MyID); ok {
		conn := conns.(Conns)
		json := append(r, []byte("\r\n--\r\n")...)
		conn.ResponseConn.Write(json)
	}
}
