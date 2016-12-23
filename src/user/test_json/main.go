package main

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"strconv"
)

func main() {
	js, err := simplejson.NewJson([]byte(`{
		"test": {
			"string_array": ["asdf", "ghjk", "zxcv"],
			"array": [1, "2", 3],
			"arraywithsubs": [{"subkeyone": 1},
			{"subkeytwo": 2, "subkeythree": 3}],
			"int": 10,
			"float": 5.150,
			"bignum": 9223372036854775807,
			"string": "simplejson",
			"bool": true
		}
	}`))
	if err != nil {
		panic("json format error")
	}
	s, err := js.Get("test").Get("string").String()
	if err != nil {
		fmt.Println("decode error: get int failed!")
		return
	}
	fmt.Println(s)

	// 检查key是否存在
	_, ok := js.CheckGet("missing_key")
	if ok {
		fmt.Println("key missing_key exists")
	} else {
		fmt.Println("key missing_key not exists")
	}

	arr, err := js.Get("test").Get("array").Array()
	if err != nil {
		fmt.Println("decode error: get array failed!")
		return
	}

	for _, v := range arr {
		var iv int
		switch v.(type) {
		case float64:
			iv = int(v.(float64))
			fmt.Println(iv)
		case string:
			iv, _ = strconv.Atoi(v.(string))
			fmt.Println(iv)
		}
	}

	//
	htmljs, err1 := simplejson.NewJson([]byte(`{}`))
	if err1 != nil {
		panic(err)
		return
	}

	htmljs.Set("appname", "qq")

	app_name, err2 := htmljs.Get("appname").String()
	if err2 != nil {
		panic(err2)
		return
	}
	fmt.Println("app_name",app_name)


}
