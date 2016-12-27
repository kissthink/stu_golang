package main

import "encoding/json"
import "fmt"


type Response2 struct {
	Page   int      `json:"page"`
	Fruits []string `json:"fruits"`
}

func main() {
	resList := []*Response2{}
	for i := 0; i < 2; i++ {
		resList = append(resList, &Response2{
			Page:   i,
			Fruits: []string{"apple", "peach", "pear"},
		})
	}
	jsonBuf, e := json.Marshal(&resList)
	fmt.Println(string(jsonBuf), e)
}
