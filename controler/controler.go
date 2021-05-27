package controler

import (
	"bytes"
	"fmt"
	"io"
	apiv1 "litrocket/api/v1"
	"litrocket/common"
	. "litrocket/common"
	"litrocket/router"
	"litrocket/utils/dataencry"
	"litrocket/utils/handlelog"
	"net"
	"strconv"
	"strings"
	"time"
)

// Strorage Url from JSON data.
type Url struct {
	Url string
}

// HandleConn Handle Client's Connection.
func HandleConn() {
	var (
		err  error    // error
		conn net.Conn // connection with client
	)

	// Loop waiting to be connected
	for {
		// Accept Client's Connection.
		if conn, err = common.RequestListener.Accept(); err != nil {
			handlelog.Handlelog("WARNING", "HandleConn Accept"+err.Error())
			continue
		}
		handlelog.Handlelog("INFO", conn.RemoteAddr().String()+" Connect ok with RequestServer")

		// We need to prevent someone from creating a large number of malicious connections.
		// When the number of connections created by an addr in 180 seconds is more than 30,
		// all the connections are closed directly

		// The connection is not malicious.
		// the new  goroutine processes the subsequent data of this connection.
		// the main goroutine continue blocking, waiting to be connected.
		go HandleData(conn)
	}
}

// HandleData handle signin and signup, after signin every single user have a goroutine,
// and every single goroutine's lifetime is user signin to signout
func HandleData(conn net.Conn) {
	var (
		n        int           // Read n bytes from client
		id       common.UserID // User ID
		signinok bool          // True SignIn and False SignUp
		err      error         // strorage error
		conns    common.Conns
		url      Url
		index    int
		last     = make([]byte, 1024*1024*3)
	)

	// defer.
	defer conn.Close()

	// Handle signin or signup.
	if id, signinok = apiv1.SignInOrUp(conn); !signinok {
		return
	}

	// Client connect to others servers.
	if err = ConnectWithServer(conn, &conns); err != nil {
		return
	}

	// Store all connections to Map.
	common.AllUsers.Store(id, conns)

	//* 开启心跳协程,每一个客户端都有一个
	go HeartBeat(id)

	last = []byte("") //没有的话,解析json会报错,有无效字符'\x00'

	// Handle data.
	for {
		buf := make([]byte, 65536)

		// Waiting client's data.
		//* 并不是每一Read都能得到准确的JSON,因为是TCP,客户端发的快可能一次Read获得多个JSON,发的慢可能一次一个JSON.
		//* 所以规定发来的JSON以\r\n--\r\n为结尾,这样就能分开了.
		if n, err = conn.Read(buf); err == io.EOF {
			// 连接断开关闭协程
			return
		}

		last = append(last, buf[0:n]...) //把整个buf[:n]追加到last

	HANDLEJSON:
		// 找不到\r\n则进行下次循环
		if index = bytes.Index(last, []byte("\r\n--\r\n")); index == -1 {
			continue
		}

		// 找出正确的JSON
		Json := bytes.Split(last, []byte("\r\n--\r\n"))

		// 在last中去掉这个json
		last = last[index+6:]

		// Parse JSON
		if err = dataencry.Unmarshal(Json[0], &url); err != nil {
			continue
		}

		// Router.
		router.Run(url.Url, Json[0])

		goto HANDLEJSON //如果last还有完整的JSON数据,继续解析,直到找不到完整的JSON数据,继续下一次Read.
	}
}

// Client Connect to Server.
func ConnectWithServer(mainconn net.Conn, conns *common.Conns) error {
	var err error

	// Get Remote IP
	i := strings.Index(mainconn.RemoteAddr().String(), ":")
	ip := mainconn.RemoteAddr().String()[:i]

	// RequestServer
	conns.RequestConn = mainconn

	// ResponseServer
	if conns.ResponseConn, err = ConnectServer(i, ip, common.ResponseListener); err != nil {
		return err
	}
	handlelog.Handlelog("INFO", "Connect ok with ResponseServer")

	// ChatServer
	if conns.ChatConn, err = ConnectServer(i, ip, common.ChatListener); err != nil {
		return err
	}
	handlelog.Handlelog("INFO", "Connect ok with ChatServer")

	// FileControl
	if conns.FileControlConn, err = ConnectServer(i, ip, common.FileControlListener); err != nil {
		return err
	}
	handlelog.Handlelog("INFO", "Connect ok with FileControl")

	// FileServer
	if conns.FileConn, err = ConnectServer(i, ip, common.FileListener); err != nil {
		return err
	}
	handlelog.Handlelog("INFO", "Connect ok with FileServer")

	// VideoServer
	if conns.VideoConn, err = ConnectServer(i, ip, common.VideoListener); err != nil {
		return err
	}
	handlelog.Handlelog("INFO", "Connect ok with VideoServer")

	return nil
}

func ConnectServer(i int, ip string, listener net.Listener) (net.Conn, error) {
	var (
		conn net.Conn
		err  error
	)

	for {
		// Accept.
		if conn, err = listener.Accept(); err != nil {
			return nil, err
		}
		// Compare ip string.
		if conn.RemoteAddr().String()[:i] == ip {
			if _, err = conn.Write([]byte("1")); err != nil {
				return nil, err
			}
			break
		}
		if _, err = conn.Write([]byte("0")); err != nil {
			return nil, err
		}
		conn.Close()
	}

	return conn, nil
}

func HeartBeat(Id common.UserID) {
	buf := make([]byte, 5)
	for {
		time.Sleep(time.Second * 3) // 3 second.
		fmt.Println("holad...")
		if conns, ok := common.AllUsers.Load(UserID(Id)); ok {
			conn := conns.(common.Conns)
			conn.RequestConn.Write([]byte("isonline")) // * Write不会发生错误总是成功的
			if _, err := conn.ResponseConn.Read(buf); err != nil {
				common.AllUsers.Delete(UserID(Id))
				handlelog.Handlelog("INFO", strconv.Itoa(int(Id))+"Leave")
				return
			}
		}
	}
}
