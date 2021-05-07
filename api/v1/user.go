package apiv1

import (
	. "litrocket/common"
	"litrocket/model"
	"litrocket/utils/dataencry"
)

// JSON Structure.
type User struct {
	Url  string
	Id   UserID
	Info int
	Mess string
}

// Return Structure User.
func GetUserInfo(json []byte) {
	var (
		user   User
		result struct {
			Url  string
			User model.User
		}
	)

	if err := dataencry.Unmarshal(json, &user); err != nil {
		return
	}

	mine, exist := model.SearchByID(user.Id)
	result.Url = user.Url
	result.User = mine

	if exist {
		if conns, ok := AllUsers.Load(user.Id); ok {
			conn := conns.(Conns)
			b, _ := dataencry.Marshal(result)
			b = append(b, []byte("\r\n--\r\n")...)
			conn.ResponseConn.Write(b)
		}
	}
}

// Edit User's Image.
func EditUserImage(json []byte) {
	var user User
	var result struct {
		Url  string
		Code int
	}

	if err := dataencry.Unmarshal(json, &user); err != nil {
		return
	}

	result.Url = user.Url
	result.Code = model.UploadUserImage(user.Mess, user.Id)

	if conns, ok := AllUsers.Load(user.Id); ok {
		conn := conns.(Conns)
		b, _ := dataencry.Marshal(result)
		conn.ResponseConn.Write(b)
	}
}

func EditUserAge(json []byte) {
	var user User
	var result struct {
		Url  string
		Code int
	}

	if err := dataencry.Unmarshal(json, &user); err != nil {
		return
	}

	result.Code = model.EditUserAge(user.Id, user.Info)
	result.Url = user.Url

	if conns, ok := AllUsers.Load(user.Id); ok {
		conn := conns.(Conns)
		b, _ := dataencry.Marshal(result)
		conn.ResponseConn.Write(b)
	}
}

func EditUserSex(json []byte) {
	var user User
	var result struct {
		Url  string
		Code int
	}

	if err := dataencry.Unmarshal(json, &user); err != nil {
		return
	}

	result.Code = model.EditUserSex(user.Id, user.Info)
	result.Url = user.Url

	if conns, ok := AllUsers.Load(user.Id); ok {
		conn := conns.(Conns)
		b, _ := dataencry.Marshal(result)
		conn.ResponseConn.Write(b)
	}
}

func EditUserTips(json []byte) {
	var user User
	var result struct {
		Url  string
		Code int
	}

	if err := dataencry.Unmarshal(json, &user); err != nil {
		return
	}

	result.Code = model.EditUserTips(user.Id, user.Mess)
	result.Url = user.Url

	if conns, ok := AllUsers.Load(user.Id); ok {
		conn := conns.(Conns)
		b, _ := dataencry.Marshal(result)
		conn.ResponseConn.Write(b)
	}
}

// 模糊搜索好友
func SearchUserByName(json []byte) {
	var user User
	var result struct {
		Users []string
	}

	if err := dataencry.Unmarshal(json, &user); err != nil {
		return
	}

	Users := model.SearchByName(user.Mess)
	for i := 0; i < len(Users); i++ {
		result.Users = append(result.Users, Users[i].UserName)
	}

	if conns, ok := AllUsers.Load(user.Id); ok {
		conn := conns.(Conns)
		b, _ := dataencry.Marshal(result)
		conn.ResponseConn.Write(b)
	}
}
