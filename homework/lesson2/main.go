package main

import (
	"fmt"
	"gb_go_best_practics/homework/lesson2/logger"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	fmt.Println("hello from lesson2")
	log.Printf("hi")

	logger.Init(os.Stdout, os.Stderr, os.Stdout, ioutil.Discard)
	logger.I.Printf("[host=%s] [uid=%d] file successfully uploaded", "srv42",
		100500)

}
