package spider

import (
	"fmt"
	"user/common"
	"testing"
	"log"
)

func Test_search(t *testing.T) {
	hc := common.NewHttpClient("https://p.xgj.me:27035")
	if hc == nil {
		log.Println("NewHttpClient..err")
		return
	}
	query_arry := []string{"facebook", "qq", "wechat", "陌陌"}

	data := make(chan string, len(query_arry))

	for _, app := range query_arry {
		go work(app, hc, data)
		//ret_list, err := search(app, hc)
		//if err != nil{
		//	return
		//}
		//fmt.Println(fmt.Sprintf("%s: %d", app,len(ret_list)))
	}
	print("just wait ...")
	num := 1
	for {
		fmt.Println("data:", <-data)
		num += 1
		if num == 4 {
			break
		}
		//time.Sleep(1e9)
		//fmt.Println("wait..")
	}
}
func work(app string, hc *common.HttpClient, data chan string) {

	ret_list, err := search(app, hc)
	if err != nil {
		data <- fmt.Sprint("%s..err", app)
		return
	}
	data <- app
	fmt.Println(fmt.Sprintf("%s: %d", app, len(ret_list)))
}
