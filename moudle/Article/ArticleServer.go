package article

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"litrocket/common"
	"litrocket/model"
	"litrocket/utils/dataencry"
	"litrocket/utils/handlelog"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const ARTICLEDIR = "tempfile/"

func InitArticle() {
	http.HandleFunc("/article", doArticle)
	http.HandleFunc("/allarticle", GetAllArticle)

	if err := http.ListenAndServe(common.ArticleAddr, nil); err != nil {
		handlelog.Handlelog("FATAL", "http ListenAndServe"+err.Error())
		return
	}
}

func doArticle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getArticle(w, r)
	case http.MethodPut:
		putArticle(w, r)
	case http.MethodDelete:
		delArticle(w, r)
	}
}

// 获取文章
func getArticle(w http.ResponseWriter, r *http.Request) {
	var result struct {
		Id      int
		Article string
	}

	// Parse From Args.
	if err := r.ParseForm(); err != nil {
		return
	}

	for k, v := range r.Form {
		fmt.Printf("[%q]=%q", k, v)
	}

	// Arg : id
	ID, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Arg : articleid
	ArtID, _ := strconv.Atoi(r.Form.Get("articleid"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Arg : securtkey
	Key := r.Form.Get("securtkey")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check Securt Key.
	h := md5.New()
	h.Write([]byte(common.AESKEY))
	if Key != base64.StdEncoding.EncodeToString(h.Sum(nil)) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Search From Db.
	arts := model.SearchArticleById(common.UserID(ArtID), common.UserID(ID))
	if len(arts) == 0 {
		w.WriteHeader(http.StatusNotFound) // 404
		return
	}

	// Open File.
	f, err := os.Open(arts[0].Content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Write Back.
	result.Id = int(arts[0].ID)
	content, _ := io.ReadAll(f)
	result.Article = string(content)

	w.WriteHeader(http.StatusOK)
	b, _ := dataencry.Marshal(result)
	w.Write(b)
}

// 获取用户的所有动态
func GetAllArticle(w http.ResponseWriter, r *http.Request) {
	var Code struct {
		Code []int
		Art  []string
	}
	if err := r.ParseForm(); err != nil {
		return
	}

	ID, _ := strconv.Atoi(r.Form.Get("id"))
	art := model.GetAllArticle(common.UserID(ID))
	artlen := len(art)
	if artlen == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	Code.Code = make([]int, artlen)
	Code.Art = make([]string, artlen)
	for i := 0; i < artlen; i++ {
		Code.Code[i] = int(art[i].ID)
		f, _ := os.Open(art[i].Content)
		b, _ := ioutil.ReadAll(f)
		Code.Art[i] = string(b)
		f.Close()
	}

	b, _ := dataencry.Marshal(Code)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// 上传文章
func putArticle(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		return
	}

	ID, _ := strconv.Atoi(r.Form.Get("id"))
	Key := r.Form.Get("securtkey")

	h := md5.New()
	h.Write([]byte(common.AESKEY))
	if Key != base64.StdEncoding.EncodeToString(h.Sum(nil)) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Create Dir.
	dir := ARTICLEDIR + "Article"
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

	if r.ContentLength > 1024*1024*5 { //* 不大于5MB
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

	if ok, id := model.CreateArticle(newfile, common.UserID(ID)); ok {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.Itoa(int(id))))
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
}

// 删除文章
func delArticle(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		return
	}

	ID, _ := strconv.Atoi(r.Form.Get("id"))
	ArtID, _ := strconv.Atoi(r.Form.Get("articleid"))
	Key := r.Form.Get("securtkey")

	h := md5.New()
	h.Write([]byte(common.AESKEY))
	if Key != base64.StdEncoding.EncodeToString(h.Sum(nil)) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	id, loc := model.DelArticle(common.UserID(ArtID), common.UserID(ID))
	if strings.Contains(loc, "../") {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	os.Remove(loc)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.Itoa(id)))
}
