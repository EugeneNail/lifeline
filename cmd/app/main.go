package main

import (
	"log"
	"time"
)

func main() {
	for {
		time.Sleep(time.Millisecond * 333)
		log.Println("It works")
	}
}
