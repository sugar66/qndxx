package main

import (
	_ "embed"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strings"
)

//go:embed index.html
var indexHtml string
var cache map[string]string

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		_, _ = io.WriteString(w, indexHtml)
		return
	}
	r.URL.Path = strings.Replace(r.URL.Path[1:], ":/", "://", 1)
	if !strings.HasPrefix(r.URL.Path, "https://h5.cyol.com/special/daxuexi/") {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "错误输入")
		return
	}

	temp := strings.Split(r.URL.Path, "/")
	temp = temp[:len(temp)-1]
	path := strings.Join(temp, "/")

	var title string
	if v, ok := cache[path]; ok {
		title = v
	} else {
		res, err := http.Get(r.URL.Path)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = io.WriteString(w, "服务异常：请求失败")
			return
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = io.WriteString(w, "服务异常：请求错误")
			return
		}
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = io.WriteString(w, "服务异常：内容解析错误")
			return
		}
		title = doc.Find("title").Text()
		if strings.Contains(title, "“青年大学习”") {
			cache[path] = title
		}
	}

	path += "/images/end.jpg"
	_, _ = fmt.Fprintf(w, `<html><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1, minimum-scale=1, maximum-scale=1, user-scalable=no"><title>%s</title></head><body style="margin:0"><div style="width:100vw;height:100vh;background-image: url(%s);background-size: 100%% 100%%;"></div></body></html>`, title, path)
}

func main() {
	cache = make(map[string]string)
	http.HandleFunc("/", index)
	_ = http.ListenAndServe(":8090", nil)
}
