package main

import (
	"github.com/bitly/go-simplejson"
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"log"
	"testing"
	"os"
	"strings"
	"net/http"
	"user/dev"
	"user/common"
)

func main() {
	htmljs, err1 := simplejson.NewJson([]byte(`{}`))
	if err1 != nil {
		panic(err1)
		return
	}

	htmljs.Set("name", "qq")

	app_name, err2 := htmljs.Get("name").String()
	if err2 != nil {
		panic(err2)
		return
	}
	fmt.Println("app_name",app_name)

	htmljs.SetPath([]string{"foo", "bar"}, "baz")


	var lst liststring.ListString
	lst.Append("hello")
	fmt.Println("%v (len: %d)", lst, lst.Len()) // [1] (len: 1)

	lst.Append("hello")
	fmt.Println("%v (len: %d)", lst, lst.Len()) // [1] (len: 1)



	hc := common.NewHttpClient("https://p.xgj.me:27035")
	if hc == nil {
		log.Println("NewHttpClient..err")
		return
	}

	ret_list, err := search("facebook", hc)
	if err != nil {
		return
	}
	fmt.Println("facebook:", len(ret_list))
}


func Test_search(t *testing.T) {
	hc := common.NewHttpClient("https://p.xgj.me:27035")
	if hc == nil {
		log.Println("NewHttpClient..err")
		return
	}
	query_arry := []string{"facebook", "qq", "wechat", "陌陌"}

	data := make(chan string, len(query_arry))

	for _, app := range query_arry {
		go work(app, hc, data)
		//ret_list, err := search(app, hc)
		//if err != nil{
		//	return
		//}
		//fmt.Println(fmt.Sprintf("%s: %d", app,len(ret_list)))
	}
	print("just wait ...")
	num := 1
	for {
		fmt.Println("data:", <-data)
		num += 1
		if num == 4 {
			break
		}
		//time.Sleep(1e9)
		//fmt.Println("wait..")
	}
}
func work(app string, hc *common.HttpClient, data chan string) {

	ret_list, err := search(app, hc)
	if err != nil {
		data <- fmt.Sprint("%s..err", app)
		return
	}
	data <- app
	fmt.Println(fmt.Sprintf("%s: %d", app, len(ret_list)))
}

/////////////////////////
//查询页面接口：
//例如	:
//		search("facebook") int {}
//return:
//		查询结果数量 & img_src_list
func search(query_app string, hc *common.HttpClient) (query_app_slice []string, err error) {

	if len(query_app) <= 0 {
		return query_app_slice, nil
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
		return query_app_slice, nil
	}

	// Create and fill the document, defer res.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		panic(err)
		return query_app_slice, nil
	}
	f, err1 := os.OpenFile(fmt.Sprintf("./test_%s.txt", query_app), os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err1 != nil {
		panic(err1)
		return query_app_slice, nil
	}
	defer f.Close()

	sum := 0
	doc.Find(".card-content").Each(func(i int, contentSelection *goquery.Selection) {
		app_deatil_url,_ := contentSelection.ChildrenFiltered(".card-click-target").Attr("href")
		if !strings.HasPrefix(app_deatil_url, "https:") {
			app_deatil_url = "https://play.google.com" + app_deatil_url
		}
		//&hl=zh_CN
		// if !strings.HasSuffix(title, "&hl=zh_CN") {
		// 	title = title + "&hl=zh-CN"
		// }
		log.Println("第", i, ":", query_app,  ":", app_deatil_url)
		//he
		// query_app_slice = append(query_app_slice, title)
		// f.WriteString(title)
		// f.WriteString("\n")
		sum += 1
		if sum >= 100 {
			return
		}
	})
	topicsSelection := doc.Find(".card .apps")
	print(topicsSelection.Length())

	fmt.Print("sum:", sum, len(query_app_slice))
	return query_app_slice[0:len(query_app_slice)], nil
}