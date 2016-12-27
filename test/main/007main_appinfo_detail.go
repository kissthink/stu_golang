package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
	"user/common"
)

type AppInfo struct {
	Developer    string `json:"developer"`
	Detail_url   string `json:"detail_url"`
	Icon         string `json:"icon"`
	Name         string `json:"name"`
	Package_name string `json:"package_name"`
	Review_stars string `json:"review_stars"`
}

func main() {
	// TestJson(nil)
	//hc := common.NewHttpClient("https://p.xgj.me:27035")
	//if hc == nil {
	//	log.Println("NewHttpClient..err")
	//	return
	//}

	startWebSer()
	//search
	//applist, err := search("facebook", hc)
	//if err != nil {
	//	panic(err)
	//	return
	//}
	//for index, app := range applist {
	//	jsonBuf, err := json.Marshal(app)
	//	if err != nil {
	//		panic(err)
	//		return
	//	}
	//	fmt.Println("第", index, "个结果是:", string(jsonBuf))
	//}
}

// url_base := "https://play.google.com/store/search?q=facebook"
func AppSearch(w http.ResponseWriter, r *http.Request) {
	hc := common.NewHttpClient("https://p.xgj.me:27035")
	if hc == nil {
		log.Println("NewHttpClient..err")
		return
	}


	r.ParseForm() //解析url传递的参数，对于POST则解析响应包的主体（request body）
	log.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		// t, _ := template.ParseFiles("login.gtpl")
		//t.Execute(w, nil)
		//for k, v := range r.Form {
		//	log.Println("key:", k)
		//	log.Println("val:", strings.Join(v, ""))
		//}
		q := r.Form.Get("q")
		log.Println(q)
		//applist, err := search("facebook", hc)
		applist, err := search(q, hc)
		if err != nil {
			panic(err)
			return
		}
		jsonBuf,err := json.Marshal(applist)
		if err != nil{
			panic(err)
			return
		}
		log.Println(string(jsonBuf))
		log.Println("\n\n\n")
		//for index, app := range applist {
		//	jsonBuf, err := json.Marshal(app)
		//	if err != nil {
		//		panic(err)
		//		return
		//	}
		//	fmt.Println("第", index, "个结果是:", string(jsonBuf))
		//}
		//fmt.Fprintf(w, "Hello search app!") //这个写入到w的是输出到客户端的
		fmt.Fprintf(w, string(jsonBuf))
	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
	fmt.Print("over")
}
func startWebSer() {
	http.HandleFunc("/search", AppSearch)    //设置search  app
	err := http.ListenAndServe("0.0.0.0:8000", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	log.Println("ser...ok")
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
				Package_name: app_package_name,
			})
			sum += 1
		})
	log.Println("sum:", sum)
	//jsonBuf, e := json.Marshal(&resList)
	//fmt.Println(string(jsonBuf), e)
	//return jsonBuf, nil
	return resList, nil
}
