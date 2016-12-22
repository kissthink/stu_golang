package main

import "fmt"
import "time"

func fibonacci(c, quit chan int) {

	x, y := 1, 1

	for {
		fmt.Println("select..bf:", x, y)
		select {

		case c <- x:
			x, y = y, x+y
			fmt.Println("x:",x,"y:",y)

		case <-quit:
			fmt.Println("quit")
			return

		}
		fmt.Println("select..af:", x, y)
	}
}

func show(c, quit chan int) {

	for i := 0; i < 10; i++ {

		fmt.Println(<-c)
	}

	quit <- 0
}

func main() {
	//
	seconds := 2
	wait := time.Duration(seconds) * time.Second

	data := make(chan int)
	leave := make(chan int)

	go show(data, leave)
	go fibonacci(data, leave)

	for {

		// time.Sleep(100)
		time.Sleep(wait)
	}



}
