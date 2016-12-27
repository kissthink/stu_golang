package spider

import (
	"github.com/PuerkitoBio/goquery"
	"os"
	"fmt"
	"strings"
	"log"
	"user/common"
	"net/http"
)
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
