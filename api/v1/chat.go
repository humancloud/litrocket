package apiv1

import (
	. "litrocket/common"
	"litrocket/model"
	"litrocket/utils/dataencry"
	"net"
)

//ChatMess strorage message
type ChatMess struct {
	Url        string
	MessFormat int    // 0 word. 1 file.
	MessDest   int    // 0 GroupMess，1 PersonalMess.
	SrcID      UserID // Sender.
	SrcName    string
	DestID     UserID // Receiver.
	DestName   string
	Strmess    string // String Messge.
	MsgDate    string // Message Date.
}

// ChatServe
func ChatServe(json []byte) {
	var (
		err      error
		chatmess ChatMess
	)

	// Parse JSON
	if err = dataencry.Unmarshal(json, &chatmess); err != nil {
		return
	}

	// Group Message
	if chatmess.MessDest == 0 {
		groupmess(json, &chatmess)
		return
	}

	// Personal Message
	personalmess(json, &chatmess)
}

// Group Message.
func groupmess(buf []byte, chatmess *ChatMess) {
	var (
		i          int64
		groupusers = make([]UserID, 200)
	)

	num := model.SearchGroupUser(groupusers, chatmess.DestID)

	for i = 0; i < num; i++ {
		chatmess.DestID = UserID(groupusers[i])
		if chatmess.DestID == chatmess.SrcID {
			continue
		}
		personalmess(buf, chatmess)
	}
}

// Personal Message.
func personalmess(buf []byte, chatmess *ChatMess) {
	destconn, r := IsOnLine(chatmess.DestID)

	switch r {
	case 0:
		// User OnLine.
		json := append(buf, []byte("\r\n--\r\n")...)
		destconn.Write(json)
	case 1:
		// If Dest User Is OffLine. Save to Db.
		str := string(buf)
		model.SaveMess(
			&str,
			chatmess.SrcID,
			chatmess.DestID,
			chatmess.MessDest,
			chatmess.MessFormat,
		)
		// Image Message.
		// if chatmess.MessFormat == 1 {
		// 	// Create Dir, dir's name is "tempdir + id"
		// 	dir := tempdir + strconv.Itoa(int(chatmess.SrcID))
		// 	if _, err := os.Stat(dir); err != nil && !os.IsExist(err) {
		// 		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		// 			handlelog.Handlelog("WARNING", "personalmess+mkdir"+err.Error())
		// 			return
		// 		}
		// 	}
		// 	// Create new file to stroage the image.
		// 	// If Dest is not online, save the image's name to StrMess.
		// 	file := dir + "/" + time.Now().Format("2006-01-02-15:04:05") + strconv.Itoa(rand.Intn(100))
		// 	if f, err = os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
		// 		handlelog.Handlelog("WARNING", "personalmess+os.open"+err.Error())
		// 		return
		// 	}

		// 	io.WriteString(f, chatmess.Strmess)
		// 	f.Close()

		// 	chatmess.Strmess = file
		// }

	default:
		return
	}
}

//IsOnLine: Dest User is Exist and Online ?
// RETURN : 0, User exists and Online;
//          1, User exists and Offline;
//         -1, User not exists.
func IsOnLine(DestID UserID) (net.Conn, int) {
	// User exists ?
	_, f := model.SearchByID(DestID)

	// User is Online ?
	conn, ok := AllUsers.Load(DestID)
	if f && ok {
		// interface{} to struct.
		conns := conn.(Conns)

		return conns.ChatConn, 0
	} else if f && !ok {
		return nil, 1
	}

	return nil, -1
}
