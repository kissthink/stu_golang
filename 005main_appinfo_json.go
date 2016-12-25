package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	//"github.com/bitly/go-simplejson"
	"log"
	"net/http"
	//"strconv"
	"strings"
	"user/common"
	"user/dev"
)

type R struct {
	Page   int                   `json:"page"`
	Fruits liststring.ListString `json:"fruits"`
}

type AppInfo struct {
	Name         string `json:"name"`
	Icon         string `json:"icon"`
	Developer    string `json:"developer"`
	Review_stars string `json:"review_stars"`
	Detail_url   string `json:"detail_url"`
}

func main() {
	// TestJson(nil)
	hc := common.NewHttpClient("https://p.xgj.me:27035")
	if hc == nil {
		log.Println("NewHttpClient..err")
		return
	}

	//search
	applist, err := search("facebook", hc)
	if err != nil {
		panic(err)
		return
	}
	//jsonBuf,err := json.Marshal(applist)
	//if err != nil{
	//	panic(err)
	//	return
	//}
	//fmt.Println(jsonBuf)

	for index, app := range applist {
		fmt.Println(app)
		jsonBuf, err := json.Marshal(app)
		if err != nil {
			panic(err)
			return
		}
		fmt.Println("第", index, "个结果是:", string(jsonBuf))
	}
}

/////////////////////////
//查询页面接口：
//例如	:
//		search("facebook") int {}
//return:
//		查询结果数量 & img_src_list
//func search(query_app string, hc *common.HttpClient) ([]byte, error) {
func search(query_app string, hc *common.HttpClient) ([]*AppInfo, error) {

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
	resList := []*AppInfo{}
	doc.Find(".card-content").Each(
		func(i int, contentSelection *goquery.Selection) {
			app_deatil_url, _ := contentSelection.ChildrenFiltered(".card-click-target").Attr("href")
			if !strings.HasPrefix(app_deatil_url, "https:") {
				app_deatil_url = "https://play.google.com" + app_deatil_url
			}

			app_icon, _ := contentSelection.Find(".cover-image").Attr("src")
			if !strings.HasPrefix(app_icon, "https:") {
				app_icon = "https:" + app_icon
			}

			app_name, _ := contentSelection.Find(".title").Attr("title")
			app_developer, _ := contentSelection.Find(".subtitle").Attr("title")
			app_package_arr := strings.Split(app_deatil_url, "?id=")

			app_package_name := ""
			if len(app_package_arr) == 2 {
				app_package_name = app_package_arr[1]
			} else {
				app_package_name = "not found id"
			}
			app_stars, _ := contentSelection.Find(".tiny-star").Attr("aria-label")

			resList = append(resList, &AppInfo{
				Name:         app_name,
				Icon:         app_icon,
				Developer:    app_developer,
				Review_stars: app_stars,
				Detail_url:   app_deatil_url,
			})

			log.Println("name:", app_name)
			log.Println("dev:", app_developer)
			log.Println("icon:", app_icon)
			log.Println("detaiurl:", app_deatil_url)
			log.Println("pakage_name:", app_package_name)
			log.Println("\n\n")
			sum += 1
		})
	log.Println("sum:", sum)
	//jsonBuf, e := json.Marshal(&resList)
	//fmt.Println(string(jsonBuf), e)
	//return jsonBuf, nil
	return resList, nil
}
