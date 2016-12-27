package main

import (
	"github.com/bitly/go-simplejson"
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"log"
	"strings"
	"net/http"
	"user/common"
	"strconv"
)

func main() {
	hc := common.NewHttpClient("https://p.xgj.me:27035")
	if hc == nil {
		log.Println("NewHttpClient..err")
		return
	}

	//search
	appListJs, err := simplejson.NewJson([]byte(`{}`))
	if err != nil {
		panic(err)
		return
	}
	appListJs, err = search("facebook", hc)
	if err != nil {
		panic(err)
		return
	}

	num,_ := appListJs.Get("sum").Int()

	for i := 0; i<num; i++{
		key := strconv.Itoa(i)
		v, _ := appListJs.Get(key).String()
		fmt.Printf("%03d...:%s", i, v)
		fmt.Println()
	}
}
/////////////////////////
//查询页面接口：
//例如	:
//		search("facebook") int {}
//return:
//		查询结果数量 & img_src_list
func search(query_app string, hc *common.HttpClient) (*simplejson.Json, error) {

	if len(query_app) <= 0 {
		return nil, nil
	}

	// url_base := "https://play.google.com/store/search?q=facebook"
	url_base := "https://play.google.com/store/search"
	url := fmt.Sprintf("%s?q=%s", url_base, query_app)
	log.Println(url)
	req, e := http.NewRequest(
		"GET",
		url,
		nil,
	)

	resp, e := hc.Do(req)
	if e != nil {
		panic(e)
		return nil, e
	}

	// Create and fill the document, defer res.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		panic(err)
		return nil, nil
	}
	sum := 0
	appListJs, err := simplejson.NewJson([]byte(`{}`))
	if err != nil {
		panic(err)
		return nil, err
	}
	doc.Find(".card-content").Each(func(i int, contentSelection *goquery.Selection) {
		app_deatil_url,_ := contentSelection.ChildrenFiltered(".card-click-target").Attr("href")
		if !strings.HasPrefix(app_deatil_url, "https:") {
			app_deatil_url = "https://play.google.com" + app_deatil_url
		}
		//&hl=zh_CN
		// if !strings.HasSuffix(title, "&hl=zh_CN") {
		// 	title = title + "&hl=zh-CN"
		// }
		//log.Println("第", i, ":", query_app,  ":", app_deatil_url)
		appListJs.Set(strconv.Itoa(i), app_deatil_url)
		sum += 1
		//if sum >= 100 {
		//	return
		//}
	})
	log.Println("sum:", sum)
	appListJs.Set("sum", sum)

	return appListJs,nil
}

