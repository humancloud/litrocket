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

		result struct {
			Url    string
			Code   int
			Friend model.User
		}

		mess struct {
			Url    string
			Friend model.User
		}
	)

	if err = dataencry.Unmarshal(json, &friend); err != nil {
		return
	}

	// 查询用户ID
	user, exist := model.SearchUser(friend.Notes)
	if !exist {
		return
	}

	result.Url = "add/friresult"
	result.Friend, result.Code = model.AddFriend(friend.SrcID, UserID(user.ID))

	// 被加的用户在线,发送消息
	if val, ok := AllUsers.Load(UserID(user.ID)); ok { //! 大坑, 如果不把ID转为UserID类型,就会检测出不在线,查询时不仅KEY的值要一样,而且KEY的类型也要和存这个键值对的时候一样
		mess.Url = friend.Url
		mess.Friend, _ = model.SearchByID(friend.SrcID)
		conns := val.(Conns)
		r, _ := dataencry.Marshal(mess)
		b := append(r, []byte("\r\n--\r\n")...)
		conns.ResponseConn.Write(b)
	}

	// 返回加好友的结果
	if val, ok := AllUsers.Load(friend.SrcID); ok {
		conns := val.(Conns)
		b, _ := dataencry.Marshal(result)
		r := append(b, []byte("\r\n--\r\n")...)
		conns.ResponseConn.Write(r)
	}
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
		b := append(json, []byte("\r\n--\r\n")...)
		conn.ResponseConn.Write(b)
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
