package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"webgp/common"
	"webgp/dev"
	"strconv"
	"github.com/PuerkitoBio/goquery"
)

type AppInfo struct {
	Developer    string `json:"developer"`
	Detail_url   string `json:"detail_url"`
	Icon         string `json:"icon"`
	Name         string `json:"name"`
	Package_name string `json:"package_name"`
	Review_stars string `json:"review_stars"`
}
type DevLink struct {
	Name string `json:"name"`
	Href string `json:"href"`
}

type AppInfoDetail struct {
	Apk_type                 string    `json:"apk_type"`
	Content_rating           dev.LStr  `json:"content_rating"`
	Category                 string    `json:"category"`
	Description              string    `json:"description"`
	Description_short        string    `json:"description_short"`
	Description_translation  string    `json:"description_translation"`
	Developer                string    `json:"developer"`
	Developer_id             string    `json:"developer_id"`
	Developer_link           []DevLink `json:"developer_link"`
	Name                     string    `json:"name"`
	Img                      dev.LStr  `json:"img"`
	Icon                     string    `json:"icon"`
	Install_count1           int       `json:"install_count_1"`
	Install_count2           int       `json:"install_count_2"`
	Interactive_elements     dev.LStr  `json:"interactive_elements"`
	In_app_products          string    `json:"in_app_products"`
	Price                    string    `json:"price"`
	Review_starts            float64   `json:"review_starts"`
	Review_count             int       `json:"review_count"`
	Published                bool      `json:"published"`
	Published_date           string    `json:"published_date"`
	Tube                     string    `json:"tube"`
	Tube_id                  string    `json:"tube_id"`
	Update                   string    `json:"update"`
	Version_description      string  `json:"version_description"`
	Version_size             string    `json:"version_size"`
	Version_current_version  string    `json:"version_current_version"`
	Version_requires_android string    `json:"version_requires_android"`
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
	r.ParseForm()                                     //解析url传递的参数，对于POST则解析响应包的主体（request body）
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
	log.Println("startWebSer...ok")
	http.HandleFunc("/search", AppSearch)           //设置search  app
	http.HandleFunc("/details", AppDetail)          //设置search  app
	err := http.ListenAndServe(":8000", nil) //设置监听的端口
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
	u, err := url.Parse("https://play.google.com/store/search")
	if err != nil {
		log.Fatal(err)
	}
	//u.Scheme = "https"
	//u.Host = "play.google.com"
	//u.Path = "/store/search"
	q := u.Query()
	q.Set("q", query_app)
	u.RawQuery = q.Encode()
	fmt.Println(u)

	req, e := http.NewRequest(
		"GET",
		u.String(),
		nil,
	)
	resp, e := hc.Do(req)
	if e != nil {
		log.Println(e)
		return nil, e
	}

	// Create and fill the document, defer res.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	sum := 0
	resList := []*AppInfo{}
	doc.Find(".card-content").Each(
		func(i int, contentSelection *goquery.Selection) {
			app_deatil_url, bool_detail_ur := contentSelection.ChildrenFiltered(".card-click-target").Attr("href")
			if !bool_detail_ur {
				log.Println(i, "app_detail_url...not exist")
				return
			}

			app_icon, bool_app_icon := contentSelection.Find(".cover-image").Attr("src")
			if !bool_app_icon {
				log.Println(i, "app_icon...not exist")
				return
			}

			app_name, bool_app_name := contentSelection.Find(".title").Attr("title")
			if !bool_app_name {
				log.Println(i, "app_name...not exist'")
				return
			}
			app_developer, bool_app_developer := contentSelection.Find(".subtitle").Attr("title")
			if !bool_app_developer {
				log.Println(i, "app_developer...not exist")
				return
			}
			app_package_arr := strings.Split(app_deatil_url, "?id=")

			app_package_name := ""
			if len(app_package_arr) == 2 {
				app_package_name = app_package_arr[1]
			} else {
				app_package_name = ""
			}
			app_stars, bool_app_starts := contentSelection.Find(".tiny-star").Attr("aria-label")
			if !bool_app_starts {
				log.Println(i, "app_starts not exist'")
				app_stars = "0"
			}

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

func detail(id string, hr string) (*AppInfoDetail, error) {
	log.Println("detail--->", "id:", id, "hr:", hr)
	u, err := url.Parse("https://play.google.com/store/apps/details")
	if err != nil {
		log.Fatal(err)
	}
	//u.Scheme = "https"
	//u.Host = "play.google.com"
	//u.Path = "/store/apps/details"
	q := u.Query()

	if len(id) > 0 {
		q.Set("id", id)
	}
	if len(hr) > 0 {
		q.Set("hr", hr)
	}
	u.RawQuery = q.Encode()
	fmt.Println(u)

	req, e := http.NewRequest(
		"GET",
		u.String(),
		nil,
	)

	hc := common.NewHttpClient("https://p.xgj.me:27035")
	if hc == nil {
		log.Println("NewHttpClient..err")
		return nil, nil
	}

	resp, e := hc.Do(req)
	if e != nil {
		log.Println(e)
		return nil, e
	}

	// Create and fill the document, defer res.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	appDetail := AppInfoDetail{}

	appDetail.Name = strings.Replace(doc.Find(".id-app-title").Text(), " ", "", -1)
	appDetail.Icon, _ = doc.Find(".cover-container .cover-image").Attr("src")

	//img
	doc.Find(".thumbnails img.screenshot").Each(func(i int, contentSelection *goquery.Selection) {
		img_url, _ := contentSelection.Attr("src")
		appDetail.Img.Append(img_url)
	})

	category, bool_detail_catagory := doc.Find(".document-subtitle.category").Attr("href")
	if !bool_detail_catagory {
		appDetail.Category = ""
		log.Println("can't find appDetail.Category")
		return nil, nil
	} else {
		appDetail.Category = category
	}
	//apk_type
	appDetail.Apk_type = "0"

	if len(appDetail.Category) > 0 {
		category_arr := strings.Split(appDetail.Category, "category/")
		if len(category_arr) >= 1 {
			appDetail.Category = category_arr[1]
		}
	}

	//content_rating
	doc.Find("[itemprop=contentRating]").Each(func(i int, contentSelection *goquery.Selection) {
		content_rating := strings.Replace(contentSelection.Text(), " ", "", -1)
		log.Println(i, content_rating)
		appDetail.Content_rating.Append(content_rating)
	})

	//describ
	js_length := doc.Find("div[itemprop=description] [jsname]").Length()
	if js_length > 0 {
		description, err_description := doc.Find("div[itemprop=description] [jsname]").Eq(0).Html()
		if err_description != nil {
			log.Println("describtion err:", err_description)
			appDetail.Description = ""
		} else {
			log.Println("desc:", description)
			appDetail.Description = strings.Replace(description, "<p></p>", "", -1)
		}

		//translation
		if js_length == 2 {
			description_translation, err_description_translation := doc.Find("div[itemprop=description] [jsname]").Eq(1).Html()
			if err_description_translation != nil {
				log.Println("describtion err:", err_description_translation)
			} else {
				appDetail.Description_translation = description_translation
			}
		}
	} else {
		description, err_description := doc.Find("div[itemprop=description]").Html()
		if err_description != nil {
			log.Println("describtion err:", err_description)
			appDetail.Description = ""
		} else {
			log.Println("desc:", description)
			appDetail.Description = strings.Replace(description, "<p></p>", "", -1)
		}
	}

	description_short, bool_detail_description_short := doc.Find("meta[name=description]").Attr("content")
	if !bool_detail_description_short {
		log.Println("description_short not exist")
		appDetail.Description_short = ""
	} else {
		appDetail.Description_short = description_short
	}

	//price
	price, bool_detail_price := doc.Find("[itemprop=price]").Attr("content")
	if !bool_detail_price {
		log.Println("price not exist")
		price = "0"
	}
	appDetail.Price = price

	//developer
	appDetail.Developer = strings.Replace(doc.Find("div[itemprop=author] [itemprop=name]").Text(), " ", "", -1)

	//devulr
	devurl, bool_detail_dev := doc.Find("div[itemprop=author] [itemprop=url]").Attr("content")
	if !bool_detail_dev {
		log.Println("devurl..not exist")
	} else {
		devurl_arr := strings.Split(devurl, "id=")
		if len(devurl_arr) >= 1 {
			appDetail.Developer_id = devurl_arr[1]
		}
	}
	other_metadata := doc.Find(".metadata")
	other_metadata.Find(".dev-link").Each(func(i int, contentSelection *goquery.Selection) {
		app_devlink := DevLink{}
		app_devlink.Href, _ = contentSelection.Attr("href")
		app_devlink.Name = strings.Replace(contentSelection.Text(), " ", "", -1)
		appDetail.Developer_link = append(appDetail.Developer_link, app_devlink)
	})

	//install
	install := strings.Replace(doc.Find("[itemprop=numDownloads]").Text(), " ", "", -1)
	install_arr := strings.Split(install, "-")
	if len(install_arr) == 2 {
		appDetail.Install_count1, _ = strconv.Atoi(strings.Replace(install_arr[0], ",", "", -1))
		appDetail.Install_count2, _ = strconv.Atoi(strings.Replace(install_arr[1], ",", "", -1))
	}

	//interactive_elements
	doc.Find(".meta-info").Each(func(i int, contentSelection *goquery.Selection) {
		title := strings.Replace(contentSelection.Find(".title").Text(), " ", "", -1)
		if len(title) > 0 {
			appDetail.Interactive_elements.Append(title)
		}
		//log.Println("title:",title)
		//in_app_products
		if title == "In-app Products" {
			appDetail.In_app_products = contentSelection.Find(".content").Text()
		} else if title == "Interactive Elements" {
			//bug
		}
	})

	//published
	appDetail.Published_date = strings.Replace(doc.Find("[itemprop=datePublished]").Text(), " ", "", -1)
	if len(appDetail.Published_date) > 0 {
		appDetail.Published = true
		appDetail.Update = appDetail.Published_date
	} else {
		appDetail.Published = false
		//需要沟通, time名，模块提供时间
		appDetail.Update = ""
	}


	//review_
	review_stars, err_review_stars := strconv.ParseFloat(doc.Find(".rating-box .score").Text(), 32)
	if err_review_stars != nil {
		log.Println("review starts", err_review_stars.Error())
		appDetail.Review_starts = 0
	} else {
		appDetail.Review_starts = review_stars
	}

	review_count, err_review_count := strconv.Atoi(strings.Replace(doc.Find(".rating-box .reviews-num").Text(), ",", "", -1))
	if err_review_count != nil {
		log.Println("review count", err_review_count.Error())
		appDetail.Review_count = 0
	} else {
		appDetail.Review_count = review_count
	}
	//tube
	doc.Find(".thumbnails span.play-action-container").Each(func(i int, contentSelection *goquery.Selection) {
		tube, bool_tube := contentSelection.Attr("data-video-url")
		if !bool_tube {
			log.Println("data-video-url not exist")
		} else {
			tube_url, err_tube_url := url.Parse(tube)
			if err_tube_url != nil {
				log.Fatal(err)
			}
			appDetail.Tube = tube_url.String()
			if len(tube_url.Path) > 0 {
				appDetail.Tube_id = strings.Replace(tube_url.Path, "/embed/", "", 1)
			}
		}
	})

	//version
	appDetail.Version_size = strings.Replace(doc.Find("[itemprop=fileSize]").Text(), " ", "", -1)
	appDetail.Version_current_version = strings.Replace(doc.Find("[itemprop=softwareVersion]").Text(), " ", "", -1)
	appDetail.Version_requires_android = strings.Replace(doc.Find("[itemprop=operatingSystems]").Text(), " ", "", -1)
	appDetail.Version_description = doc.Find(".recent-change").Text()

	return &appDetail, nil
}
