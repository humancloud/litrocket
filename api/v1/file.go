package apiv1

import (
	"fmt"
	"io"
	"math/rand"
	"strings"

	. "litrocket/common"
	"litrocket/model"
	"litrocket/utils/dataencry"
	"litrocket/utils/errmsg"
	"litrocket/utils/handlelog"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	tempfile = "tempfile/"
)

//Filetran
type Filetran struct {
	Url      string
	Filename string
	SrcName  string
	Size     int64
	SrcID    UserID
	DestID   UserID
}

// Upload Group File
func UploadGroupFile(json []byte) {
	var (
		file Filetran
		err  error
	)

	if err = dataencry.Unmarshal(json, &file); err != nil {
		return
	}

	// New Dir, dir's name is "group + groupid"
	dir := tempfile + "group" + strconv.Itoa(int(file.DestID))
	if _, err = os.Stat(dir); err != nil && !os.IsExist(err) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			handlelog.Handlelog("WARNING", "uploadgroupfile+mkdir"+err.Error())
			return
		}
	}

	// New File, File's Name Is "tempfile/filename-time-rand number"
	newfile := dir + "/" + file.Filename + time.Now().Format("2006-01-02-15:04:02") + strconv.Itoa(rand.Intn(100))
	f, err := os.Create(newfile)
	if err != nil {
		handlelog.Handlelog("WARNING", "groupfile"+"addgroupfile"+"os.Create"+err.Error())
		return
	}

	//* 创建新协程来接收文件,不卡主协程,不然文件大的话,时间会很长.
	go func() {
		var r struct {
			Url  string
			Code int
		}

		// Receive File.
		if conns, ok := AllUsers.Load(file.SrcID); ok {
			conn := conns.(Conns)
			r.Code = receive(f, conn.FileConn, file.Size)
			r.Url = file.Url
			result, _ := dataencry.Marshal(r)
			b := append(result, []byte("\r\n--\r\n")...)
			conn.FileControlConn.Write(b)
		}

		// Save Record to DataBase.
		if r.Code == errmsg.OK_SUCCESS {
			now := time.Now().Format("2006-01-02 15:04:02")
			record := model.File{
				Filetime:     now,
				FileSize:     file.Size,
				Filename:     file.Filename,
				SrcName:      file.SrcName,
				FileLocation: newfile,
				SrcID:        file.SrcID,
				DestID:       file.DestID,
				FileType:     model.GROUPFILE,
				FileState:    model.STROAGE}
			model.UploadGroupFile(&record)
		}

		f.Close() //* 最后关文件
	}()
}

// View Group Files
func ViewGroupFiles(json []byte) {
	var (
		lens   int
		file   Filetran
		result struct {
			Url   string
			Code  int // -1 failed. errmsg.OK_SUCCESS success.
			Files []struct {
				FileName string
				FileSize int64
				FileTime string
				SrcName  string
			}
		}
	)

	if err := dataencry.Unmarshal(json, &file); err != nil {
		return
	}

	// View
	files := model.ViewGroupFiles(file.DestID)
	lens = len(files)

	if lens == 0 {
		result.Code = -1
	} else {
		result.Code = errmsg.OK_SUCCESS
	}

	// 分配内存
	result.Files = make([]struct {
		FileName string
		FileSize int64
		FileTime string
		SrcName  string
	}, lens)
	for i := 0; i < lens; i++ {
		result.Files[i].FileName = files[i].Filename
		result.Files[i].FileSize = files[i].FileSize
		result.Files[i].FileTime = files[i].Filetime
		result.Files[i].SrcName = files[i].SrcName
	}

	result.Url = file.Url
	b, _ := dataencry.Marshal(result)
	// Get Conn
	if conns, ok := AllUsers.Load(file.SrcID); ok {
		conn := conns.(Conns)
		r := append(b, []byte("\r\n--\r\n")...)
		conn.FileControlConn.Write(r)
	}
}

// Delete Group File.
func DeleteGroupFile(json []byte) {
	var (
		file     Filetran
		location string
		result   struct {
			Url  string
			Code int
		}
	)

	if err := dataencry.Unmarshal(json, &file); err != nil {
		return
	}

	// 删除数据库记录
	result.Url = file.Url
	result.Code, location = model.DeleteGroupFile(file.Filename, file.SrcID, file.DestID)

	// 删除文件
	if location != "" && -1 == strings.Index(location, "../") {
		if location[0:9] == tempfile {
			os.Remove(location)
		}
	}

	// 发送结果
	if conns, ok := AllUsers.Load(file.SrcID); ok {
		conn := conns.(Conns)
		r, _ := dataencry.Marshal(result)
		b := append(r, []byte("\r\n--\r\n")...)
		conn.FileControlConn.Write(b)
	}
}

// Download Group File
func DownloadGroupFile(json []byte) {
	var (
		fileconn net.Conn
		file     Filetran
	)

	if err := dataencry.Unmarshal(json, &file); err != nil {
		return
	}

	filelocation, filelen := model.GetGroupFileLocation(file.Filename, file.DestID)

	if conns, ok := AllUsers.Load(file.SrcID); ok {
		conn := conns.(Conns)
		fileconn = conn.FileConn
	}

	if filelocation != "" && -1 == strings.Index(filelocation, "../") {
		if f, err := os.Open(filelocation); err == nil {
			go sendfile(f, fileconn, filelen)
		}
	}
}

// Send Personal File
func SendPersonalFile(json []byte) {
	var (
		file Filetran
		r    struct {
			Url  string
			Code int
		}
	)

	if err := dataencry.Unmarshal(json, &file); err != nil {
		return
	}

	// New Dir
	dir := tempfile + "User" + strconv.Itoa(int(file.SrcID))
	if _, err := os.Stat(dir); err != nil && !os.IsExist(err) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			handlelog.Handlelog("WARNING", "sendpersonalfile+mkdir"+err.Error())
			return
		}
	}

	newfile := dir + "/" + file.Filename + time.Now().Format("2006-01-02-15:04:02") + strconv.Itoa(rand.Intn(100))
	f, err := os.Create(newfile)
	if err != nil {
		handlelog.Handlelog("WARNING", "sendpersonalfile"+" os.Create:"+err.Error())
		return
	}

	go func() {
		// Receive File.
		if conns, ok := AllUsers.Load(file.SrcID); ok {
			conn := conns.(Conns)
			r.Url = file.Url
			r.Code = receive(f, conn.FileConn, file.Size)
			result, _ := dataencry.Marshal(r)
			b := append(result, []byte("\r\n--\r\n")...)
			conn.FileControlConn.Write(b)
		}

		// Send Mess.
		file.Url = "file/come"
		b, _ := dataencry.Marshal(file)
		res := append(b, []byte("\r\n--\r\n")...)

		if conns, ok := AllUsers.Load(file.DestID); ok {
			conn := conns.(Conns)
			conn.FileControlConn.Write(res)
		} else {
			str := string(res)
			model.SaveMess(&str, file.SrcID, file.DestID, 1, 1)
		}

		// Save to DB
		now := time.Now().Format("2006-01-02 15:04:02")
		record := model.File{
			Filetime:     now,
			FileSize:     file.Size,
			Filename:     file.Filename,
			SrcName:      file.SrcName,
			FileLocation: newfile,
			SrcID:        file.SrcID,
			DestID:       file.DestID,
			FileType:     model.PERSONFILE,
			FileState:    model.NOSEND}

		model.PersonFile(&record)

		f.Close()
	}()
}

// Recv Personal File
// User select download personalfile
func RecvPersonalFile(json []byte) {
	var (
		file Filetran
	)

	if err := dataencry.Unmarshal(json, &file); err != nil {
		return
	}

	filelocation, size := model.GetFileLocation(file.Filename, file.SrcID, file.DestID)
	f, err := os.Open(filelocation)
	if err != nil {
		handlelog.Handlelog("WARNING", "recvpersonalfile"+err.Error())
		return
	}

	if conns, ok := AllUsers.Load(file.DestID); ok {
		conn := conns.(Conns)
		go sendfile(f, conn.FileConn, size)
	}
}

// Receive File From Client.
func receive(f *os.File, conn net.Conn, filelen int64) int {
	var (
		n        int
		err      error
		filesize int64
	)
	buf := make([]byte, 65536)

	// Read and Write
	for filesize = 0; filesize < filelen; filesize += int64(n) {

		// Read From Client.
		if n, err = conn.Read(buf); err != nil {
			fmt.Println("read", err.Error())
			return errmsg.ERR_RECEIVE_ERROR
		}

		// Write to File
		if _, err = f.Write(buf[:n]); err != nil {
			fmt.Println("write", err.Error())
			return errmsg.ERR_RECEIVE_ERROR
		}
	}

	// md5
	return errmsg.OK_SUCCESS
}

// Send File to Dest
func sendfile(f *os.File, conn net.Conn, filelen int64) int {
	var (
		size int64
		n    int
		err  error
		buf  = make([]byte, 65536)
	)

	//Write...
	for size = 0; size < filelen; size += int64(n) {
		// Read From File.
		if n, err = f.Read(buf); err != nil {
			if err == io.EOF {
				break
			}
			return errmsg.ERR_SEND_ERROR
		}

		// Write to Dest.
		if _, err = conn.Write(buf[:n]); err != nil {
			return errmsg.ERR_SEND_ERROR
		}

		// Sleep 125ms 限速
		time.Sleep(time.Millisecond * 125)
	}

	f.Close()

	return errmsg.OK_SUCCESS
}
