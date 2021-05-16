package apiv1

import (
	. "litrocket/common"
	"litrocket/model"
	"litrocket/utils/dataencry"
)

type Dict struct {
	Url     string
	UserID  UserID
	Chinese string
	English string
}

// 根据中文模糊搜索单词并返回匹配的单词
// todo 结果特别多的情况下(超过1MB的数据,三千条记录才200K),分批返回,估计不会特别大,因为单词数量还较小.
func PersonalChiDict(json []byte) {
	var (
		dict   Dict
		result struct {
			Url   string
			Dicts []model.Dict
		}
	)

	if err := dataencry.Unmarshal(json, &dict); err != nil {
		return
	}

	result.Url = dict.Url
	result.Dicts = model.ChiLikeSearch(dict.Chinese, dict.UserID)

	if conns, ok := AllUsers.Load(dict.UserID); ok { // ! 存到Map时是 UserID 取时就算是common.UserID都不行哦
		conn := conns.(Conns)
		b, _ := dataencry.Marshal(result)
		r := append(b, []byte("\r\n--\r\n")...)
		conn.ResponseConn.Write(r)
	}
}

func PersonalEngDict(json []byte) {
	var (
		dict   Dict
		result struct {
			Url   string
			Dicts model.Dict
		}
	)

	if err := dataencry.Unmarshal(json, &dict); err != nil {
		return
	}

	result.Url = dict.Url
	result.Dicts = model.EngSearch(dict.English, dict.UserID)

	if conns, ok := AllUsers.Load(dict.UserID); ok {
		conn := conns.(Conns)
		b, _ := dataencry.Marshal(result)
		r := append(b, []byte("\r\n--\r\n")...)
		conn.ResponseConn.Write(r)
	}
}

func PushWord(json []byte) {
	var (
		dict   Dict
		result struct {
			Url  string
			Code int
		}
	)
	if err := dataencry.Unmarshal(json, &dict); err != nil {
		return
	}

	result.Url = dict.Url
	result.Code = model.PushWord(dict.Chinese, dict.English, dict.UserID)

	if conns, ok := AllUsers.Load(dict.UserID); ok {
		conn := conns.(Conns)
		b, _ := dataencry.Marshal(result)
		r := append(b, []byte("\r\n--\r\n")...)
		conn.ResponseConn.Write(r)
	}
}
