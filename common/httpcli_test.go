package main

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestHttpCliGet(t *testing.T) {
	hc := NewHttpClient("https://p.xgj.me:27035")
	url := "https://www.google.com"
	req, e := http.NewRequest(
		"Get",
		url,
		nil,
	)
	resp, e := hc.Do(req)
	if e != nil {
		t.Fatal(e)
	}
	defer resp.Body.Close()
	buf, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		t.Fatal(e)
	}
	t.Log(string(buf))
}

