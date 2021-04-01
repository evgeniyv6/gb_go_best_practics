package main

import (
	"log"
	"time"
)

func main() {
	var files chan string
	// 	for bad -->>
	//fmock := []string{
	//	path.Join(curDir, "files/1.txt"),
	//	path.Join(curDir, "files/2.txt"),
	//	path.Join(curDir, "files/3.txt"), // bad
	//	path.Join(curDir, "files/5.txt"),
	//}
	//
	//for _, ff := range fmock {
	//	files <- ff
	//}
	//
	//completed := computeHashes(files)
	//
	//<-completed
	//log.Println("done")

	// good --->
	done := make(chan interface{})
	terminated := goodHash(done, files)
	go func() {
		time.Sleep(1*time.Second)
		log.Println("canceling computeHashesTerminatable goroutine")
		close(done)
	}()

	<-terminated
	log.Println("done")

}
