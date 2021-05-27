package apiv1

import (
	"fmt"
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
			Json []string
		}
	)

	if err := dataencry.Unmarshal(json, &m); err != nil {
		return
	}

	result.Url = m.Url

	mes := model.SendOffLineMess(m.Id)
	size := len(mes)
	result.Json = make([]string, size)
	for i := 0; i < size; i++ {
		result.Json[i] = mes[i].Message
	}

	b, _ := dataencry.Marshal(result)

	fmt.Println(string(b))

	if conns, ok := AllUsers.Load(UserID(m.Id)); ok {
		conn := conns.(Conns)
		conn.ResponseConn.Write(b)
		conn.ResponseConn.Write([]byte("\r\n--\r\n"))
	}
}
