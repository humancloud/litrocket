package apiv1

import (
	"fmt"
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
			Code int
		}
	)

	if err := dataencry.Unmarshal(json, &group); err != nil {
		return
	}

	// 查询名称是否重复
	if ok := model.SearchGroupExist(group.GroupName); ok {
		result.Code = -1
		goto WRITE
	}

	G = model.GroupInfo{GroupName: group.GroupName, GroupRootID: group.RootID, GroupUserNum: 1}
	id = model.CreateGroupInfo(G)
	g = model.Group{GroupID: id, UserRole: 0, UserID: group.RootID, UserState: 0}
	model.CreateGroup(g)

	result.Code = 0

WRITE:
	r, _ := dataencry.Marshal(result)
	if conns, ok := AllUsers.Load(group.RootID); ok {
		conn := conns.(Conns)
		conn.RequestConn.Write(r)
	}
}

func AddGroup(json []byte) {
	var (
		group Group
		err   error
	)

	if err = dataencry.Unmarshal(json, &group); err != nil {
		return
	}

	// 查询群聊id
	id := model.SearchGroupIdByName(group.GroupName)

	// 直接添加到数据库
	model.AddGroup(group.MyID, id)

	// 获取管理员ID,管理员在线则发送消息,不在线则进入数据库
	//if conns, ok := AllUsers.Load(group.MyID); ok {
	//	conn := conns.(Conns)
	//	conn.ResponseConn.Write(json)
	//}
}

func AddGroupOk(json []byte) {
	var (
		group  Group
		err    error
		result GroupResult
	)

	if err = dataencry.Unmarshal(json, &group); err != nil {
		return
	}

	i := model.AddGroupOk(group.MyID, group.GroupID)
	result.Code = i
	result.GroupID = group.GroupID
	result.Url = group.Url

	r, _ := dataencry.Marshal(result)

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
		group Group
		err   error
	)

	if err = dataencry.Unmarshal(json, &group); err != nil {
		return
	}

	model.DelGroup(group.MyID, group.GroupID)
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
	r, _ := dataencry.Marshal(result)

	fmt.Println(string(r))

	if conns, ok := AllUsers.Load(group.MyID); ok {
		conn := conns.(Conns)
		conn.RequestConn.Write(r)
	}
}
