package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func computeHashes(files <-chan string) <-chan interface{} {
	completed := make(chan interface{})

	go func() {
		defer close(completed)
		for f := range files {
			h,err := compute(f)
			if err != nil{
				log.Printf("err %v hash for file %q", err, f)
				continue
			}
			log.Printf("hash for file %q is %d", f,h)
		}
	}()

	return completed
}

func compute(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", errors.New(fmt.Sprintf("read file %q err %v", f, err))
	}

	h:= sha256.New()
	if _, err:=io.Copy(h, f); err != nil {
		return "", errors.New(fmt.Sprintf("compute hash err %s",err))
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
