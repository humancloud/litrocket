package apiv1

import (
	"litrocket/common"
	. "litrocket/common"
	"litrocket/model"
	"litrocket/utils/dataencry"
)

type OffLineMess struct {
	Url string
	Id  common.UserID
}

func SendOffLineMess(json []byte) {
	var (
		m      OffLineMess
		result struct {
			Url  string
			Date []string
			Json []string
		}
	)

	if err := dataencry.Unmarshal(json, &m); err != nil {
		return
	}

	result.Url = m.Url

	mes := model.SendOffLineMess(m.Id)
	size := len(mes)
	result.Date = make([]string, size)
	result.Json = make([]string, size)
	for i := 0; i < size; i++ {
		result.Date[i] = mes[i].MsgDate
		result.Json[i] = mes[i].Message
	}

	b, _ := dataencry.Marshal(result)

	if conns, ok := AllUsers.Load(UserID(m.Id)); ok {
		conn := conns.(Conns)
		conn.RequestConn.Write(b)
		conn.RequestConn.Write([]byte("\r\n--\r\n"))
	}
}
