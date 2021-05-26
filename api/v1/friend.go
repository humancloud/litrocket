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
// 某用户发起好友推荐
// 将用户根据其好友关系组成一个图
// 以此用户为顶点,寻找一度好友(我的好友)和每一个一度好友的二度好友(好友的好友)
// 如果有两个或多个一度好友有同样的二度好友(除自己),即推荐这个同样的二度好友.  也就是共同好友(多数都是这样)

// 查询数据库一度好友,添加至邻接表
func FriendRecommand(json []byte) {
	var (
		ReComd struct {
			Url string
			Id  UserID
		}

		result struct {
			Url  string
			Name []string
		}
	)

	if err := dataencry.Unmarshal(json, &ReComd); err != nil {
		return
	}

	result.Url = ReComd.Url

	//
}
