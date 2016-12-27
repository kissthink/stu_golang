package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bitly/go-simplejson"
	"log"
	"net/http"
	"strconv"
	"strings"
	"user/common"

	"user/dev"
)

type R struct {
	Page   int      `json:"page"`
	//Fruits []string `json:"fruits"`
	Fruits liststring.ListString `json:"fruits"`
}

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

	num, _ := appListJs.Get("sum").Int()

	for i := 0; i < num; i++ {
		key := strconv.Itoa(i)
		app := appListJs.Get(key)
		println(app)
		app_name, _ := app.Get("name").String()
		app_icon, _ := app.Get("icon").String()
		app_developer, _ := app.Get("developer").String()
		app_stars, _ := app.Get("review_stars").String()

		fmt.Printf("%03d...:%s", i, app_name)
		fmt.Printf("%03d...:%s", i, app_icon)
		fmt.Printf("%03d...:%s", i, app_developer)
		fmt.Printf("%03d...:%s", i, app_stars)
		fmt.Println()
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
	doc.Find(".card-content").Each(
		func(i int, contentSelection *goquery.Selection) {
		APPNEW:
			appJs, err1 := simplejson.NewJson([]byte(`{}`))
			if err1 != nil {
				panic(err1)
				//continue
				goto APPNEW
			}
			app_deatil_url, _ := contentSelection.ChildrenFiltered(".card-click-target").Attr("href")
			if !strings.HasPrefix(app_deatil_url, "https:") {
				app_deatil_url = "https://play.google.com" + app_deatil_url
			}

			appJs.Set("detail_url", app_deatil_url)
			app_icon, _ := contentSelection.Find(".cover-image").Attr("src")
			if !strings.HasPrefix(app_icon, "https:") {
				app_icon = "https:" + app_icon
			}
			appJs.Set("icon", app_icon)
			//&hl=zh_CN
			// if !strings.HasSuffix(title, "&hl=zh_CN") {
			// 	title = title + "&hl=zh-CN"
			// }
			//log.Println("第", i, ":", query_app,  ":", app_deatil_url)
			app_name, _ := contentSelection.Find(".title").Attr("title")
			appJs.Set("name", app_name)

			app_developer, _ := contentSelection.Find(".subtitle").Attr("title")
			appJs.Set("developer", app_developer)

			app_package_arr := strings.Split(app_deatil_url, "?id=")
			app_package_name := ""
			if len(app_package_arr) == 2 {
				app_package_name = app_package_arr[1]
			} else {
				app_package_name = "not found id"
			}
			appJs.Set("package_name", "facebook")

			app_stars, _ := contentSelection.Find(".tiny-star").Attr("aria-label")
			appJs.Set("review_stars", app_stars)

			log.Println("name:", app_name)
			log.Println("dev:", app_developer)
			log.Println("icon:", app_icon)
			log.Println("detaiurl:", app_deatil_url)
			log.Println("pakage_name:", app_package_name)
			log.Println("\n\n")

			appListJs.Set(strconv.Itoa(i), appJs)
			sum += 1
		})
	log.Println("sum:", sum)
	appListJs.Set("sum", sum)

	return appListJs, nil
}
