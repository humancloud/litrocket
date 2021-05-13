package apiv1

import (
	"container/list"
	"fmt"
	"io/ioutil"
	. "litrocket/common"
	"litrocket/utils/dataencry"
	"litrocket/utils/handlelog"
	"net"
	"net/http"
	"strconv"
)

const (
	checkkey = "http://127.0.0.1::8090/control/get?room="
)

type Video struct {
	Url    string
	SrcID  UserID
	DestID UserID
}

type Join struct {
	User UserID
	Conn net.Conn
}

var joinuser = make(map[UserID]*list.List, 200)

// 加入共享，将加入结果发回
func JoinScreen(json []byte) {
	var (
		v      Video
		ok     bool
		err    error
		conns  interface{}
		result struct {
			Code int
		}
	)

	// Json
	if err = dataencry.Unmarshal(json, &v); err != nil {
		return
	}

	// 未开始共享
	if _, ok = joinuser[v.DestID]; !ok {
		result.Code = -1
		return
	}

	if conns, ok = AllUsers.Load(v.SrcID); ok {
		conn := conns.(Conns)
		newuser := Join{User: v.SrcID, Conn: conn.VideoConn}
		list := joinuser[v.DestID]
		list.PushBack(newuser)
		result.Code = 0
		r, _ := dataencry.Marshal(result)
		conn.VideoConn.Write(r)
	}
}

// 客户端停止拉流即可，同时在链表中去除
func QuitScreen(json []byte) {
	var (
		video Video
		err   error
	)

	if err = dataencry.Unmarshal(json, &video); err != nil {
		return
	}

	if v, ok := joinuser[video.DestID]; ok {
		for e := v.Front(); e != nil; e = e.Next() {
			s := e.Value.(Join)
			if s.User == video.SrcID {
				v.Remove(e)
				fmt.Println("Quit", joinuser)
				fmt.Println(joinuser[video.SrcID])
			}
		}
	}
}

// 在livego中创建room,开始共享
func SendScreen(json []byte) {
	var (
		video Video
		err   error
		resp  *http.Response
		buf   = make([]byte, 65536)
		Json  struct {
			Status int
			Data   string
		}
		result struct {
			Code int
		}
	)

	if err = dataencry.Unmarshal(json, &video); err != nil {
		return
	}

	joins := list.New()
	joinuser[video.SrcID] = joins

	resp, err = http.Get(checkkey + strconv.Itoa(int(video.SrcID)))
	if err != nil {
		handlelog.Handlelog("WARNING", "http.Get(), checkkey"+err.Error())
		result.Code = -1
		goto WRITE
	}

	buf, _ = ioutil.ReadAll(resp.Body)
	dataencry.Unmarshal(buf, &Json)
	if Json.Status != 200 {
		result.Code = -1
		goto WRITE
	}

WRITE:
	r, _ := dataencry.Marshal(result)

	if conns, ok := AllUsers.Load(video.SrcID); ok {
		conn := conns.(Conns)
		conn.VideoConn.Write(r)
	}
}

// 结束共享，删除map数据即可，不再推流后livego的流会自动停止
func EndScreen(json []byte) {
	var (
		video Video
		err   error
	)

	if err = dataencry.Unmarshal(json, &video); err != nil {
		return
	}

	if _, ok := joinuser[video.SrcID]; ok {
		// 销毁链表
		if v, ok := joinuser[video.SrcID]; ok {
			for e := v.Front(); e != nil; e = e.Next() {
				v.Remove(e)
			}
		}
		delete(joinuser, video.SrcID)
	}
}

func SeaVideo(json []byte) {

}
