package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"os"
	"strings"
)

func main() {
	url_base := "http://studygolang.com/topics"
	url := url_base

	f, err1 := os.OpenFile("studygolang.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err1 != nil {
		panic(err1)
		return
	}

	defer f.Close()

	//var buf []string

	for i := 0; i < 147; i++ {
		doc, err := goquery.NewDocument(url)
		if err != nil {
			log.Fatal(err)
		}

		doc.Find(".topics .topic").Each(func(i int, contentSelection *goquery.Selection) {
			title := contentSelection.Find(".title a").Text()
			log.Println("第", i+1, "个帖子的标题：", title)
			buf := []string {"第", strconv.Itoa(i+1), "个帖子的标题：", title}
			tmp := strings.Join(buf, "")
			f.WriteString(tmp)
			f.WriteString("\n")

		})

		url = url_base + "?p=" + strconv.Itoa(i)
		fmt.Println(url)
		f.WriteString(url)
		f.WriteString("\n")
		f.WriteString("\n")

	}

}
