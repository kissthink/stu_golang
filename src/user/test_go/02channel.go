package main

import (
	"fmt"
	"time"
)

func show(c chan int) {
	// seconds := 1
	// wait := time.Duration(seconds) * time.Second

	for {
		fmt.Print("af read from c:")
		fmt.Println()
		data := <-c
		fmt.Print("bf read from c:", data)
		fmt.Println()
		if 1 == data {
			fmt.Print("receive ")
			fmt.Println()
		}
	}
}

func main() {
	//
	seconds := 2
	wait := time.Duration(seconds) * time.Second

	c := make(chan int)

	go show(c)

	for {
		time.Sleep(wait)
		c <- 1
		fmt.Print("send ")
		fmt.Println()
	}
}
