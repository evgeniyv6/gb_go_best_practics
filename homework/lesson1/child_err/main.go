package main

import (
	"log"
	"path"
)

func main() {
	done := make(chan interface{})
	defer close(done)

	files := []string{
		path.Join(curDir, "files/1.txt"),
		path.Join(curDir, "files/2.txt"),
		path.Join(curDir, "files/3.txt"), // bad
		path.Join(curDir, "files/5.txt"),
	}

	for res := range badHasher(done, files...) {
		log.Printf("bad hash res %q", res)
	}

	log.Println("------")

	for res := range myFilesHasher(done, files...) {
		if res.err != nil {
			log.Printf("get hash err %v", res.err)
			break
		}
		log.Printf("file hash %q", res.hash)
	}

}
