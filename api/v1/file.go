package apiv1

import (
	"io"
	"math/rand"

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

// uploadgroupfile
func UploadGroupFile(json []byte) {
	var (
		file Filetran
		err  error
		r    struct {
			Url  string
			Code int
		}
	)

	if err = dataencry.Unmarshal(json, &file); err != nil {
		return
	}

	// 判断名字是否重复，重复不可上传
	// New Dir, dir's name is "group + groupid"
	dir := tempfile + "group" + strconv.Itoa(int(file.DestID))
	if _, err = os.Stat(dir); err != nil && !os.IsExist(err) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			handlelog.Handlelog("WARNING", "uploadgroupfile+mkdir"+err.Error())
			return
		}
	}

	// New File.
	newfile := dir + "/" + file.Filename + time.Now().Format("2006-01-02-15:04:02") + strconv.Itoa(rand.Intn(100))
	f, err := os.Create(newfile)
	if err != nil {
		handlelog.Handlelog("WARNING", "groupfile"+"addgroupfile"+"os.Create"+err.Error())
		return
	}
	defer f.Close()

	// Receive File.
	if conns, ok := AllUsers.Load(file.SrcID); ok {
		conn := conns.(Conns)
		r.Code = receive(f, conn.FileConn, file.Size)
		r.Url = file.Url
		result, _ := dataencry.Marshal(r)
		conn.ResponseConn.Write(result)
		conn.ResponseConn.Write([]byte("\r\n--\r\n"))
	}

	// Save Record to DataBase.
	if r.Code == errmsg.OK_SUCCESS {
		now := time.Now().Format("2006-01-02 15:04:02")
		file := model.File{
			Filetime:     now,
			FileSize:     file.Size,
			Filename:     file.Filename,
			SrcName:      file.SrcName,
			FileLocation: newfile,
			SrcID:        file.SrcID,
			DestID:       file.DestID,
			FileType:     model.GROUPFILE,
			FileState:    model.STROAGE}
		model.UploadGroupFile(&file)
	}
}

// viewgroupfiles
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
		conn.ResponseConn.Write(b)
		conn.ResponseConn.Write([]byte("\r\n--\r\n"))
	}
}

// deletegroupfile
func DeleteGroupFile(json []byte) {
	var (
		file   Filetran
		result struct {
			Code int // -1 failed. errmsg.OK_SUCCESS success.
		}
	)

	if err := dataencry.Unmarshal(json, &file); err != nil {
		return
	}

	// 暂时不删除文件，只删除数据库记录
	result.Code = model.DeleteGroupFile(file.Filename, file.DestID)

	if conns, ok := AllUsers.Load(file.SrcID); ok {
		conn := conns.(Conns)
		r, _ := dataencry.Marshal(result)
		conn.FileConn.Write(r)
	}
}

// downloadgroupfile
func DownloadGroupFile(json []byte) {
	var (
		fileconn net.Conn
		file     Filetran
		result   struct {
			Code int
		}
	)

	if err := dataencry.Unmarshal(json, &file); err != nil {
		return
	}

	if conns, ok := AllUsers.Load(file.SrcID); ok {
		conn := conns.(Conns)
		fileconn = conn.FileConn
	}

	filelocation, filelen := model.GetGroupFileLocation(file.Filename, file.DestID)

	if filelocation != "" {
		if f, err := os.Open(filelocation); err == nil {
			sendfile(f, fileconn, filelen)
		}

		result.Code = errmsg.ERR_FILE_NO_EXIST
		b, _ := dataencry.Marshal(result)
		fileconn.Write(b)
	}
}

// sendpersonalfile
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
	defer f.Close()

	// Receive File.
	if conns, ok := AllUsers.Load(file.SrcID); ok {
		conn := conns.(Conns)
		r.Url = file.Url
		r.Code = receive(f, conn.FileConn, file.Size)
		result, _ := dataencry.Marshal(r)

		conn.ResponseConn.Write(result)
		conn.ResponseConn.Write([]byte("\r\n--\r\n"))
	}

	// Send Mess
	if conns, ok := AllUsers.Load(file.DestID); ok {
		conn := conns.(Conns)
		conn.FileConn.Write(json)
		conn.FileConn.Write([]byte("\r\n--\r\n"))
		return
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
}

// recvpersonalfile
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
	defer f.Close()

	if conns, ok := AllUsers.Load(file.DestID); ok {
		conn := conns.(Conns)
		sendfile(f, conn.FileConn, size)
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
			return errmsg.ERR_RECEIVE_ERROR
		}

		// Write to File
		if _, err = f.Write(buf[:n]); err != nil {
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
		buf  = make([]byte, 10*1024*1024)
	)

	//Write...
	for size = 0; size != filelen; size += int64(n) {
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
	}

	return errmsg.OK_SUCCESS
}
