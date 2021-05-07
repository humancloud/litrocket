package apiv1

import (
	. "litrocket/common"
	"litrocket/model"
	"litrocket/utils/dataencry"
)

type Friend struct {
	Url    string
	SrcID  UserID
	DestID UserID
	Notes  string
}

type FriResult struct {
	Url      string
	Code     int
	FriendID UserID
}

// Add a friend to table "friend", friend's state is "waiting friend agree".
func AddFriend(json []byte) {
	var (
		err    error
		friend Friend
	)

	if err = dataencry.Unmarshal(json, &friend); err != nil {
		return
	}

	// 查询用户ID
	user, exist := model.SearchUser(friend.Notes)
	if !exist {
		return
	}

	model.AddFriend(friend.SrcID, UserID(user.ID))

	// 被加的用户在线,发送消息
	if val, ok := AllUsers.Load(friend.DestID); ok {
		conns := val.(Conns)
		conns.ResponseConn.Write(json)
	}

	// 不在线,暂时存到数据库表,登录时显示
}

// I Agree friend request.
func Agree(json []byte) {
	var (
		err    error
		friend Friend
		r      []byte
		result FriResult
	)
	if err := dataencry.Unmarshal(json, &friend); err != nil {
		return
	}

	i := model.AddFriendOk(friend.SrcID, friend.DestID)

	// 加我的用户在线,则发送通知,不在线只需更改数据库
	if val, ok := AllUsers.Load(friend.DestID); ok {
		conns := val.(Conns)

		result.Code = i
		result.FriendID = friend.DestID
		result.Url = friend.Url

		if r, err = dataencry.Marshal(result); err != nil {
			return
		}
		conns.ResponseConn.Write(r)
	}
}

// 不同意好友请求
func NoAgree(json []byte) {
	// 在数据库中删除即可,不需其他操作
}

func GetAllFriend(json []byte) {
	var (
		friend  Friend
		friends []model.Friend
		result  struct {
			Url    string
			Friend []model.Friend
		}
	)

	if err := dataencry.Unmarshal(json, &friend); err != nil {
		return
	}

	friends = model.GetAllFriend(friend.SrcID)

	result.Friend = friends
	result.Url = friend.Url
	buf, _ := dataencry.Marshal(result)
	buf = append(buf, []byte("\r\n--\r\n")...)

	if conns, ok := AllUsers.Load(friend.SrcID); ok {
		conn := conns.(Conns)
		conn.ResponseConn.Write(buf)
	}
}

func GetFriendInfo(json []byte) {
	var (
		exist  bool
		friend struct {
			Url      string
			SrcID    UserID
			FriendID UserID
		}
	)

	var (
		result struct {
			Url  string
			Info model.User
		}
	)

	if err := dataencry.Unmarshal(json, &friend); err != nil {
		return
	}

	result.Url = friend.Url
	result.Info, exist = model.SearchByID(friend.FriendID)
	if exist {
		if conns, ok := AllUsers.Load(friend.SrcID); ok {
			conn := conns.(Conns)
			buf, _ := dataencry.Marshal(result)
			buf = append(buf, []byte("\r\n--\r\n")...)
			conn.ResponseConn.Write(buf)
		}
	}

}

func DelFriend(json []byte) {
	var (
		friend Friend
	)

	if err := dataencry.Unmarshal(json, &friend); err != nil {
		return
	}

	model.DelFriend(friend.SrcID, friend.DestID)

	// 目标在线则通知,不在线不通知
	if conns, ok := AllUsers.Load(friend.DestID); ok {
		conn := conns.(Conns)
		conn.ResponseConn.Write(json)
	}
}

func EditFriendNotes(json []byte) {
	var (
		friend Friend
		result FriResult
	)

	if err := dataencry.Unmarshal(json, &friend); err != nil {
		return
	}

	i := model.EditFriendNotes(friend.SrcID, friend.DestID, friend.Notes)
	result.Code = i
	result.FriendID = friend.DestID
	result.Url = friend.Url
	r, _ := dataencry.Marshal(result)

	if conns, ok := AllUsers.Load(friend.SrcID); ok {
		conn := conns.(Conns)
		conn.ResponseConn.Write(r)
	}
}

// Friend Recommand
func FriendRecommand(json []byte) {

}
