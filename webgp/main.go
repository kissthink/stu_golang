package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"webgp/common"
	"webgp/dev"

	"flag"

	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
)

type AppBase struct {
	Developer   string `json:"developer,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Name        string `json:"name,omitempty"`
	PackageName string `json:"package_name,omitempty"`
}

type AppInfo struct {
	AppBase
	DetailUrl   string `json:"detail_url,omitempty"`
	ReviewStars string `json:"review_stars,omitempty"`
}

type DevLink struct {
	Name string `json:"name,omitempty"`
	Href string `json:"href,omitempty"`
}

type AppInfoDetail struct {
	AppBase
	ApkType                string    `json:"apk_type,omitempty"`
	ContentRating          dev.LStr  `json:"content_rating,omitempty"`
	Category               string    `json:"category,omitempty"`
	Description            string    `json:"description,omitempty"`
	DescriptionShort       string    `json:"description_short,omitempty"`
	DescriptionTranslation string    `json:"description_translation,omitempty"`
	DeveloperId            string    `json:"developer_id,omitempty"`
	DeveloperLink          []DevLink `json:"developer_link,omitempty"`
	Img                    dev.LStr  `json:"img,omitempty"`
	InstallCount1          int       `json:"install_count_1,omitempty"`
	InstallCount2          int       `json:"install_count_2,omitempty"`
	InteractiveElements    dev.LStr  `json:"interactive_elements,omitempty"`
	InAppProducts          string    `json:"in_app_products,omitempty"`
	Price                  string    `json:"price,omitempty"`
	ReviewStars            float64   `json:"review_stars,omitempty"`
	ReviewCount            int       `json:"review_count,omitempty"`
	PubliShed              bool      `json:"published,omitempty"`
	PubliShedDate          string    `json:"published_date,omitempty"`
	Tube                   string    `json:"tube,omitempty"`
	TubeId                 string    `json:"tube_id,omitempty"`
	Update                 string    `json:"update,omitempty"`
	VersionDescription     string    `json:"version_description,omitempty"`
	VersionSize            string    `json:"version_size,omitempty"`
	VersionCurrentVersion  string    `json:"version_current_version,omitempty"`
	VersionRequiresAndroid string    `json:"version_requires_android,omitempty"`
}

func main() {
	startWebSer()
}

func startWebSer() {
	//初始化命令行参数
	flag.Parse()

	glog.Info("StartWebSer...start")
	http.HandleFunc("/search", AppSearch)    //设置search  app
	http.HandleFunc("/details", AppDetail)   //设置detail  app
	err := http.ListenAndServe(":8000", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	//flag.Parse()
}

// url_base := "https://play.google.com/store/search?q=facebook"
func AppSearch(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "GET" {
		q := r.Form.Get("q")
		applist, err := search(q)
		if err != nil {
			glog.Errorln(err)
			return
		}
		jsonBuf, err := json.Marshal(applist)
		if err != nil {
			glog.Errorln(err)
			return
		}
		fmt.Fprintf(w, string(jsonBuf))
	} else {
		glog.Warning("not found")
	}
	glog.Info("app search finish ok")
	//退出时调用，确保日志写入文件中
	defer glog.Flush()
}

//detail_url: "https://play.google.com/store/apps/details?id=com.tencent.mm",
func AppDetail(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "GET" {
		id := r.Form.Get("id")
		hl := r.Form.Get("hl")
		appdetail, err := detail(id, hl)
		if err != nil {
			glog.Errorln(err)
			return
		}
		jsonBuf, err := json.Marshal(appdetail)
		if err != nil {
			glog.Errorln(err)
			return
		}
		//glog.Info(string(jsonBuf))
		//glog.Info("\n\n\n")
		fmt.Fprintf(w, string(jsonBuf))
	} else {
		glog.Warning("not found")
	}
	glog.Info("app detail finish ok")
	defer glog.Flush()
}

/////////////////////////
//查询页面接口：
//例如	:
//		search("facebook") int {}
//return:
//		查询结果数量 & img_src_list
func search(appName string) ([]*AppInfo, error) {

	hc := common.NewHttpClient("https://p.xgj.me:27035")
	if hc == nil {
		glog.Errorln("NewHttpClient..err")
		return nil, nil
	}

	if len(appName) <= 0 {
		glog.Errorln("query_app can't be empty")
		return nil, nil
	}

	u, err := url.Parse("https://play.google.com/store/search")
	if err != nil {
		glog.Errorln(err)
	}

	q := u.Query()
	q.Set("q", appName)
	u.RawQuery = q.Encode()

	req, e := http.NewRequest(
		"GET",
		u.String(),
		nil,
	)

	resp, e := hc.Do(req)
	if e != nil {
		glog.Errorln(e)
		return nil, e
	}

	// Create and fill the document, defer res.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		glog.Errorln(err)
		return nil, nil
	}

	sum := 0
	resList := []*AppInfo{}
	doc.Find(".card-content").Each(
		func(i int, contentSelection *goquery.Selection) {
			appDeatilUrl, errBool := contentSelection.ChildrenFiltered(".card-click-target").Attr("href")
			if !errBool {
				glog.Errorln(i, "app_detail_url...not exist")
				return
			}

			appIcon, errBool := contentSelection.Find(".cover-image").Attr("src")
			if !errBool {
				glog.Errorln(i, "app_icon...not exist")
				return
			}

			appName, errBool := contentSelection.Find(".title").Attr("title")
			if !errBool {
				glog.Errorln(i, "app_name...not exist'")
				return
			}
			appDeveloper, errBool := contentSelection.Find(".subtitle").Attr("title")
			if !errBool {
				glog.Errorln(i, "app_developer...not exist")
				return
			}
			appPackageArr := strings.Split(appDeatilUrl, "?id=")

			appPackageName := ""
			if len(appPackageArr) == 2 {
				appPackageName = appPackageArr[1]
			} else {
				appPackageName = ""
			}
			appStars, errBool := contentSelection.Find(".tiny-star").Attr("aria-label")
			if !errBool {
				glog.Info(i, "app_stars not exist'")
				appStars = "0"
			}

			appbase := AppBase{
				Name:        strings.TrimSpace(appName),
				Icon:        strings.TrimSpace(appIcon),
				Developer:   strings.TrimSpace(appDeveloper),
				PackageName: strings.TrimSpace(appPackageName),
			}
			resList = append(resList, &AppInfo{
				AppBase:     appbase,
				DetailUrl:   strings.TrimSpace(appDeatilUrl),
				ReviewStars: strings.TrimSpace(appStars),
			})
			sum += 1
		})
	glog.Info("sum:", sum)
	return resList, nil
}

func detail(id string, hr string) (*AppInfoDetail, error) {
	u, err := url.Parse("https://play.google.com/store/apps/details")
	if err != nil {
		log.Fatal(err)
	}

	q := u.Query()

	if len(id) > 0 {
		q.Set("id", id)
	}
	if len(hr) > 0 {
		q.Set("hr", hr)
	}

	u.RawQuery = q.Encode()
	glog.Info(u)

	req, e := http.NewRequest(
		"GET",
		u.String(),
		nil,
	)

	hc := common.NewHttpClient("https://p.xgj.me:27035")
	if hc == nil {
		glog.Errorln("NewHttpClient..err")
		return nil, nil
	}

	resp, e := hc.Do(req)
	if e != nil {
		glog.Errorln(e)
		return nil, e
	}

	// Create and fill the document, defer res.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		glog.Errorln(err)
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

	category, errBool := doc.Find(".document-subtitle.category").Attr("href")
	if !errBool {
		appDetail.Category = ""
		glog.Errorln("can't find appDetail.Category")
		return nil, nil
	} else {
		appDetail.Category = category
	}

	//apk_type
	appDetail.ApkType = "0"

	if len(appDetail.Category) > 0 {
		category_arr := strings.Split(appDetail.Category, "category/")
		if len(category_arr) >= 1 {
			appDetail.Category = category_arr[1]
		}
	}

	//content_rating
	doc.Find("[itemprop=contentRating]").Each(func(i int, contentSelection *goquery.Selection) {
		content_rating := strings.Replace(contentSelection.Text(), " ", "", -1)
		glog.Info(i, content_rating)
		appDetail.ContentRating.Append(content_rating)
	})

	//describ
	js_length := doc.Find("div[itemprop=description] [jsname]").Length()
	if js_length > 0 {
		description, err := doc.Find("div[itemprop=description] [jsname]").Eq(0).Html()
		if err != nil {
			glog.Info("describtion err:", err)
			appDetail.Description = ""
		} else {
			//glog.Info("desc:", description)
			appDetail.Description = strings.Replace(description, "<p></p>", "", -1)
		}

		//translation
		if js_length == 2 {
			description_translation, err := doc.Find("div[itemprop=description] [jsname]").Eq(1).Html()
			if err != nil {
				glog.Info("describtion err:", err)
			} else {
				appDetail.DescriptionTranslation = description_translation
			}
		}
	} else {
		description, err := doc.Find("div[itemprop=description]").Html()
		if err != nil {
			glog.Info("describtion err:", err)
			appDetail.Description = ""
		} else {
			glog.Info("desc:", description)
			appDetail.Description = strings.Replace(description, "<p></p>", "", -1)
		}
	}

	description_short, errBool := doc.Find("meta[name=description]").Attr("content")
	if !errBool {
		glog.Info("description_short not exist")
		appDetail.DescriptionShort = ""
	} else {
		appDetail.DescriptionShort = description_short
	}

	//price
	price, errBool := doc.Find("[itemprop=price]").Attr("content")
	if !errBool {
		glog.Info("price not exist")
		price = "0"
	}
	appDetail.Price = price

	//developer
	appDetail.Developer = strings.Replace(doc.Find("div[itemprop=author] [itemprop=name]").Text(), " ", "", -1)

	//devulr
	devurl, errBool := doc.Find("div[itemprop=author] [itemprop=url]").Attr("content")
	if !errBool {
		log.Println("devurl..not exist")
	} else {
		devurl_arr := strings.Split(devurl, "id=")
		if len(devurl_arr) >= 1 {
			appDetail.DeveloperId = devurl_arr[1]
		}
	}
	other_metadata := doc.Find(".metadata")
	other_metadata.Find(".dev-link").Each(func(i int, contentSelection *goquery.Selection) {
		app_devlink := DevLink{}
		app_devlink.Href, _ = contentSelection.Attr("href")
		app_devlink.Name = strings.Replace(contentSelection.Text(), " ", "", -1)
		appDetail.DeveloperLink = append(appDetail.DeveloperLink, app_devlink)
	})

	//install
	install := strings.Replace(doc.Find("[itemprop=numDownloads]").Text(), " ", "", -1)
	install_arr := strings.Split(install, "-")
	if len(install_arr) == 2 {
		appDetail.InstallCount1, _ = strconv.Atoi(strings.Replace(install_arr[0], ",", "", -1))
		appDetail.InstallCount2, _ = strconv.Atoi(strings.Replace(install_arr[1], ",", "", -1))
	}

	//interactive_elements
	doc.Find(".meta-info").Each(func(i int, contentSelection *goquery.Selection) {
		title := contentSelection.Find(".title").Text()
		if len(title) > 0 {
			title = strings.TrimSpace(title)
			//in_app_products
			if title == "In-app Products" {
				appDetail.InAppProducts = contentSelection.Find(".content").Text()
				//glog.Info("in_app_products:", appDetail.In_app_products)
			} else if title == "Interactive Elements" {
				title_tmp := contentSelection.Find(".content").Text()
				if len(title_tmp) > 0 {
					//glog.Info("todo:", title_tmp)
					appDetail.InteractiveElements = strings.Split(title_tmp, ", ")
				}
			}
		}

	})

	//published
	appDetail.PubliShedDate = strings.Replace(doc.Find("[itemprop=datePublished]").Text(), " ", "", -1)
	if len(appDetail.PubliShedDate) > 0 {
		appDetail.PubliShed = true
		appDetail.Update = appDetail.PubliShedDate
	} else {
		appDetail.PubliShed = false
		//TODO需要沟通, time名，模块提供时间
		appDetail.Update = ""
	}

	//review_
	review_stars, err_review_stars := strconv.ParseFloat(doc.Find(".rating-box .score").Text(), 32)
	if err_review_stars != nil {
		glog.Errorln("review starts", err_review_stars.Error())
		appDetail.ReviewStars = 0
	} else {
		appDetail.ReviewStars = review_stars
	}

	review_count, err := strconv.Atoi(strings.Replace(doc.Find(".rating-box .reviews-num").Text(), ",", "", -1))
	if err != nil {
		glog.Errorln("review count", err.Error())
		appDetail.ReviewCount = 0
	} else {
		//glog.Info("review count", review_count)
		//glog.Info("review count-src:", doc.Find(".rating-box .reviews-num").Text())
		appDetail.ReviewCount = review_count
	}

	//tube
	doc.Find(".thumbnails span.play-action-container").Each(func(i int, contentSelection *goquery.Selection) {
		tube, errBool := contentSelection.Attr("data-video-url")
		if !errBool {
			glog.Errorln("data-video-url not exist")
		} else {
			tube_url, err_tube_url := url.Parse(tube)
			if err_tube_url != nil {
				glog.Errorln(err_tube_url)
			}
			appDetail.Tube = tube_url.String()
			if len(tube_url.Path) > 0 {
				appDetail.TubeId = strings.Replace(tube_url.Path, "/embed/", "", 1)
			}
		}
	})

	//version
	appDetail.VersionSize = strings.Replace(doc.Find("[itemprop=fileSize]").Text(), " ", "", -1)
	appDetail.VersionCurrentVersion = strings.Replace(doc.Find("[itemprop=softwareVersion]").Text(), " ", "", -1)
	appDetail.VersionRequiresAndroid = strings.TrimSpace(doc.Find("[itemprop=operatingSystems]").Text())
	appDetail.VersionDescription = doc.Find(".recent-change").Text()

	return &appDetail, nil
}
