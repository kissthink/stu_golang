package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"testing"
)

func TestHttpCliGet(t *testing.T) {
	//hc := NewHttpClient("https://p.xgj.me:27035")
	hc := NewHttpClient("ip://192.168.1.102")
	//url := "https://play.google.com/store/apps/search?q=qq"
	url := "https://www.baidu.com"
	req, e := http.NewRequest(
		"Get",
		url,
		nil,
	)
	resp, e := hc.Do(req)
	if e != nil {
		log.Println(url)
		t.Fatal(e)
	}
	defer resp.Body.Close()
	buf, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		t.Fatal(e)
	}
	t.Log(string(buf))
}

func TestHttpCliPath(t *testing.T) {
	var key string
	var resUrl string

	key = "encrypt key"
	resUrl = "https://www.google.com"

	resUri, pErr := url.Parse(resUrl)
	if pErr != nil {
		fmt.Print(pErr)
		return
	}

	fmt.Println("resUri:", resUri)

	path := resUri.EscapedPath()
	fmt.Println("path:", path)

	rawStr := fmt.Sprintf("%s%s%s", resUri, key, path)

	fmt.Println(rawStr)
}
