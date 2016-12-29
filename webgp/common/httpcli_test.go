package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"
	"webgp/waiter"
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

func TestHttpCliGetWaiter(t *testing.T) {
	//hc := NewHttpClient("https://p.xgj.me:27035")
	var wg sync.WaitGroup
	hc := NewHttpClient("ip://192.168.1.102")
	hc.Waiter = waiter.NewBurstLimitTick(time.Second, 3)
	time.Sleep(3 * time.Second)
	b := time.Now()

	for i := 0; i < 9; i++ {
		wg.Add(1)
		go func() {
			<-hc.Waiter.GetC()
			println("i:", i, time.Now().String())
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(time.Since(b))
}
