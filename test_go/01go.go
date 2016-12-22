package main

import "fmt"
import "time"

func show() {

	seconds := 1
	wait := time.Duration(seconds) * time.Second

	for {
		fmt.Print("child 1")
		fmt.Println()
		time.Sleep(wait)
	}
}
func main() {

	seconds := 1
	wait := time.Duration(seconds) * time.Second

	go show()

	for {
		fmt.Print("parent 1")
		fmt.Println()
		time.Sleep(wait)
	}
}
