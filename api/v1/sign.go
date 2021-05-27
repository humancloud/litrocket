package apiv1

import (
	"crypto/md5"
	"encoding/base64"
	"hash"
	. "litrocket/common"
	"litrocket/model"
	"litrocket/utils/dataencry"
	"litrocket/utils/errmsg"
	"litrocket/utils/handlelog"
	"net"
)

// Sign strorage SignIn or SignUp info.
type Sign struct {
	Messtype   int    // 0 is SignUp, 1 is SignIn. 2 is forgot password. 3 is SignUp Verify Code. 4 is ForGotPassWd Verify code.
	Name       string // User name.
	Passwd     string // Password.
	Mail       string // Mail
	VerifyCode string // VerifyCode
}

// SignInOrUp
func SignInOrUp(conn net.Conn) (UserID, bool) {
	var (
		n    int
		err  error
		sign Sign // strorage signin or signup info in JSON.
	)
	buf := make([]byte, 1024)

	// Read Data From Client.
	if n, err = conn.Read(buf); err != nil {
		conn.Close()
		return -1, false
	}

	// Parse JSON data.
	if err = dataencry.Unmarshal(buf[0:n], &sign); err != nil {
		conn.Close()
		return -1, false
	}

	switch sign.Messtype {
	case 0: //SignUp
		SignUp(conn, &sign)
		return -1, false
	case 1: //SignIn
		id, code := SignIn(conn, &sign)
		if code != errmsg.OK_SUCCESS {
			return -1, false
		}

		return UserID(id), true
	case 2: // forgot password.
		ForGotPassWd(conn, &sign)
	case 3: // 注册用验证码
		SignUpVerify(sign.Mail)
		return -1, false
	case 4: // 重置密码用验证码
		ReloadPassVerify(sign.Mail)
		return -1, false
	}

	conn.Close()
	return -1, false
}

//SignUp
func SignUp(conn net.Conn, sign *Sign) {
	var (
		passwd string
		result struct {
			Code int
		}
		h         hash.Hash
		md5passwd string
	)

	// Search User by User's name
	if _, exists := model.SearchUser(sign.Name); exists {
		result.Code = errmsg.ERR_NAMEREPEAT
		goto WRITE
	}

	// Search Mail
	if ok := model.SearchMail(sign.Mail); ok {
		result.Code = errmsg.ERR_MAILREPEAT
		goto WRITE
	}

	// Check Verify Code
	if signupcode[sign.Mail] != sign.VerifyCode {
		result.Code = errmsg.ERR_CODEERROR
		goto WRITE
	}

	// Verify ok, delete this key-value.
	codemutex.Lock()
	delete(signupcode, sign.Mail)
	codemutex.Unlock()

	// Add user and info to database.
	passwd = dataencry.DecryptPasswd(sign.Passwd)

	// Password Md5 Stroage In Db.
	h = md5.New()
	h.Write([]byte(passwd))
	md5passwd = base64.StdEncoding.EncodeToString(h.Sum(nil))

	_ = model.InsertUser(sign.Name, md5passwd, sign.Mail)
	result.Code = errmsg.OK_SUCCESS
	handlelog.Handlelog("INFO", sign.Name+" SignUp success")

WRITE:
	// Write SignUp result to client
	r, _ := dataencry.Marshal(result)
	conn.Write(r)
	conn.Close()
}

//SignIn
func SignIn(conn net.Conn, sign *Sign) (uint, int) {
	var (
		id     uint
		user   model.User
		passwd string
		h      hash.Hash
	)
	var result struct {
		Code int
		ID   UserID
	}

	// Search User by User's name
	user, exists := model.SearchUser(sign.Name)
	if !exists {
		result.Code = errmsg.ERR_NOSUCHUSER
		goto WRITE
	}

	// Alealdy signin
	if v, ok := AllUsers.Load(UserID(user.ID)); v != nil || ok { // 注意:Load的参数类型一定要与Store时的key的类型一致，不然是查不到的.
		result.Code = errmsg.ERR_ALEALDYSIGNIN
		goto WRITE
	}

	// Verify password and user id.
	passwd = dataencry.DecryptPasswd(sign.Passwd)
	// Password Md5 Stroage In Db.
	h = md5.New()
	h.Write([]byte(passwd))
	passwd = base64.StdEncoding.EncodeToString(h.Sum(nil))
	if user.PassWord != passwd {
		result.Code = errmsg.ERR_WRONGPASSWD
		goto WRITE
	}

	// All signin info is right.
	result.Code = errmsg.OK_SUCCESS
	result.ID = UserID(user.ID)
	id = user.ID
	handlelog.Handlelog("INFO", sign.Name+" Login")

WRITE:
	// Write SignIn result to client.
	r, _ := dataencry.Marshal(result)
	conn.Write(r)
	return id, result.Code
}

func ForGotPassWd(conn net.Conn, sign *Sign) {
	var (
		h      hash.Hash
		passwd string
		result struct {
			Code int
		}
	)

	// Search User by User's name
	user, exists := model.SearchUser(sign.Name)
	if !exists {
		result.Code = errmsg.ERR_NOSUCHUSER
		goto WRITE
	}

	// Check Mail
	if sign.Mail != user.UserMail {
		result.Code = errmsg.ERR_MAILWRONG
		goto WRITE
	}

	// VerifyCode
	if forgotcode[user.UserMail] != sign.VerifyCode {
		result.Code = errmsg.ERR_CODEERROR
		goto WRITE
	}

	// Update DB.
	passwd = dataencry.DecryptPasswd(sign.Passwd)

	// Password Md5 Stroage In Db.
	h = md5.New()
	h.Write([]byte(passwd))
	passwd = base64.StdEncoding.EncodeToString(h.Sum(nil))
	result.Code = model.UpDatePasswd(passwd, sign.Name)

WRITE:
	r, _ := dataencry.Marshal(result)
	conn.Write(r)
	conn.Close()
}
