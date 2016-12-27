package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
	"user/common"
	"user/dev"
)

type lstring liststring.ListString

type AppInfo struct {
	Developer    string `json:"developer"`
	Detail_url   string `json:"detail_url"`
	Icon         string `json:"icon"`
	Name         string `json:"name"`
	Package_name string `json:"package_name"`
	Review_stars string `json:"review_stars"`
}
type AppInfoDetail struct {
	Apk_type		string	`json:"apk_type"`
	Content_rating		lstring	`json:"content_rating"`
	Category		string	`json:"category"`
	Description 		string	`json:"description"`
	Description_short 	string	`json:"description_short"`
	Description_ranslation	string	`json:"description_ranslation"`
	Developer		string	`json:"developer"`
	Developer_link		string	`json:"developer_link"`
	Name         		string	`json:"name"`
	Img			lstring	`json:"img"`
	Icon         		string	`json:"icon"`
	Install_count1		int	`json:"install_count_1"`
	Install_count2		int	`json:"install_count_2"`
	Interactive_elements	lstring	`json:"interactive_elements"`
	In_app_products		string	`json:"in_app_products"`
	Price			string	`json:"price"`
	Review_starts		float32	`json:"review_starts"`
	Review_count		int	`json:"review_count"`
	Published		bool	`json:"published"`
	Tube			string	`json:"tube"`
	Tube_id			string	`json:"tube_id"`
	Update			string	`json:"update"`
	Version_description	string	`json:"version_description"`
	Version_size		string	`json:"version_size"`
	Version_current_version	string	`json:"version_current_version"`
	Version_requires_android string	`json:"version_requires_android"`
}

func main() {
	startWebSer()
}

// url_base := "https://play.google.com/store/search?q=facebook"
func AppSearch(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()                    //解析url传递的参数，对于POST则解析响应包的主体（request body）
	log.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		q := r.Form.Get("q")
		log.Println(q)
		applist, err := search(q)
		if err != nil {
			panic(err)
			return
		}
		jsonBuf, err := json.Marshal(applist)
		if err != nil {
			panic(err)
			return
		}
		log.Println(string(jsonBuf))
		log.Println("\n\n\n")
		fmt.Fprintf(w, string(jsonBuf))
	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		fmt.Println("not found")
	}
	fmt.Print("over")
}

//detail_url: "https://play.google.com/store/apps/details?id=com.tencent.mm",
func AppDetail(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()                    //解析url传递的参数，对于POST则解析响应包的主体（request body）
	log.Println("detail-view:---->method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		id := r.Form.Get("id")
		hl := r.Form.Get("hl")
		log.Println("args:", id, hl)
		appdetail, err := detail(id, hl)
		if err != nil {
			panic(err)
			return
		}
		jsonBuf, err := json.Marshal(appdetail)
		if err != nil {
			panic(err)
			return
		}
		log.Println(string(jsonBuf))
		log.Println("\n\n\n")
		fmt.Fprintf(w, string(jsonBuf))
	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		fmt.Println("not found")
	}
	fmt.Print("over")

}

func startWebSer() {
	http.HandleFunc("/search", AppSearch)           //设置search  app
	http.HandleFunc("/details", AppDetail)          //设置search  app
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
func search(query_app string) ([]*AppInfo, error) {

	hc := common.NewHttpClient("https://p.xgj.me:27035")
	if hc == nil {
		log.Println("NewHttpClient..err")
		return nil, nil
	}


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
			//if !strings.HasPrefix(app_deatil_url, "https:") {
			//	app_deatil_url = "https://play.google.com" + app_deatil_url
			//}

			app_icon, _ := contentSelection.Find(".cover-image").Attr("src")
			//if !strings.HasPrefix(app_icon, "https:") {
			//	app_icon = "https:" + app_icon
			//}

			app_name, _ := contentSelection.Find(".title").Attr("title")
			app_developer, _ := contentSelection.Find(".subtitle").Attr("title")
			app_package_arr := strings.Split(app_deatil_url, "?id=")

			app_package_name := ""
			if len(app_package_arr) == 2 {
				app_package_name = app_package_arr[1]
			} else {
				app_package_name = ""
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
	return resList, nil
}

func detail(id string, hr string) ([]*AppInfoDetail, error) {
	log.Println("detail--->","id:", id, "hr:", hr)

	url_base := "https://play.google.com/store/apps/details"
	var url string
	if len(id) > 0 {
		url = fmt.Sprintf("%s%s%s", url_base,"?id=", id)
	}

	if len(hr) > 0 {
		url = fmt.Sprintf("%s%s%s", url_base,"&hr=", id)
	}

	hc := common.NewHttpClient("https://p.xgj.me:27035")
	if hc == nil {
		log.Println("NewHttpClient..err")
		return nil, nil
	}

	log.Println("details_url:",url)
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

	print(doc)

	return nil, nil
}