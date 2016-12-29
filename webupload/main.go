package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
	"os"
	"github.com/golang/glog"
	"html/template"
)

// 处理/upload 逻辑
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		fmt.Fprintf(w, "<h1>上传OK<h1>")
		f, err := os.OpenFile("./upload/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

func main() {
	glog.Info("StartWebSer...start")
	http.HandleFunc("/upload", upload)
	err := http.ListenAndServe(":9000", nil) //设置监听的端口
	if err != nil {
		glog.Error("ListenAndServe: ", err)
	}
}
