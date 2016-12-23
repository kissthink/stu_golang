package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"log"
	"work/common"
)

func main()  {
	hc := NewHttpClient("https://p.xgj.me:27035")
	url := "https://www.google.com"
	req, e := http.NewRequest(
		"Get",
		url,
		nil,
	)
	resp, e := hc.Do(req)
	if e != nil {
		print(e)
	}
	defer resp.Body.Close()
	buf, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		log.Printf(e)
	}

	fmt.Print(buf)

}