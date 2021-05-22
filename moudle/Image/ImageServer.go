package image

import (
	"crypto/md5"
	"encoding/base64"
	"io/ioutil"
	"litrocket/common"

	"litrocket/utils/dataencry"
	"litrocket/utils/handlelog"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const tempfile = "tempfile/"

// Init Image Server.
func InitImageServer() {
	http.HandleFunc("/", GetImg)
	http.HandleFunc("/image", UploadImg)

	if err := http.ListenAndServe(common.ImageAddr, nil); err != nil {
		handlelog.Handlelog("FATAL", "http ListenAndServe"+err.Error())
		return
	}
}

func GetImg(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		doGet(w, r)
	}
}

func UploadImg(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		PutImg(w, r)
	}
}

// GetImage
func doGet(w http.ResponseWriter, r *http.Request) {
	var (
		file string
		f    *os.File
		err  error
		buf  = make([]byte, 65536)
	)

	if file = r.URL.Path[1:]; file == "" {
		return
	}

	if f, err = os.Open(file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	w.WriteHeader(http.StatusOK)
	for {
		n, err := f.Read(buf)
		if n == 0 || err != nil {
			break
		}
		w.Write(buf[:n])
	}
}

// Put Image
func PutImg(w http.ResponseWriter, r *http.Request) {
	var result struct {
		Link string
	}

	if err := r.ParseForm(); err != nil {
		return
	}

	Key := r.Form.Get("securtkey")

	h := md5.New()
	h.Write([]byte(common.AESKEY))
	if Key != base64.StdEncoding.EncodeToString(h.Sum(nil)) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Create Dir.
	dir := tempfile + "Image"
	if _, err := os.Stat(dir); err != nil && !os.IsExist(err) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			handlelog.Handlelog("WARNING", "Image+mkdir"+err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// Create File.
	newfile := dir + "/" + time.Now().Format("2006-01-02-15:04:02") + strconv.Itoa(rand.Intn(100))
	f, err := os.Create(newfile)
	if err != nil {
		handlelog.Handlelog("WARNING", "image"+"addimage"+"os.Create"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	if r.ContentLength > 1024*1024*10 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write to File
	f.Write(content)

	result.Link = "http://" + "127.0.0.1" + common.ImageAddr + "/" + newfile

	b, _ := dataencry.Marshal(result)

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
