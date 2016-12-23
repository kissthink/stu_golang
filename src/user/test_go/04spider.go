package main

import (
	//"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	//"strconv"
	//"os"
	//"strings"
	//"go/doc"
	"fmt"
	//"work/common"
)

func main() {
	url_base := "https://play.google.com/store/search?q=facebook"
	url := url_base


	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	html := doc.Find("card-click-target ").Text()
	fmt.Print(html)


}

