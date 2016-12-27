package main

import (
	"fmt"
	"github.com/bitly/go-simplejson"
)




func main() {
	//
	htmljs, err1 := simplejson.NewJson([]byte(`{}`))
	if err1 != nil {
		panic(err)
		return
	}

	htmljs.Set("name", "qq")

	app_name, err2 := htmljs.Get("name").String()
	if err2 != nil {
		panic(err2)
		return
	}
	fmt.Println("app_name",app_name)

	htmljs.SetPath([]string{"foo", "bar"}, "baz")

}

