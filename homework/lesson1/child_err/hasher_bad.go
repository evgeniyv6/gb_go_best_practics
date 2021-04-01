package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func badHasher(done <-chan interface{}, files ...string) <-chan string {
	hashes := make(chan string)

	go func() {
		defer close(hashes)
		for _ , f := range files {
			hash, err := compute(f)
			if err != nil {
				log.Printf("get hash err %s", err)
				continue
			}

			select {
			case <-done:
				return
			case hashes <-hash:

			}
		}
	}()

	return hashes
}

func compute(path string) (string, error){
	f, err := os.Open(path)
	if err != nil {
		return "", errors.New(fmt.Sprintf("error: %s, path: %q", err, path))
	}

	h:= sha256.New()
	if _, err := io.Copy(h,f); err != nil {
		return "", errors.New(fmt.Sprintf("compute hash err %s",err))
	}

	return fmt.Sprintf("%x",h.Sum(nil)), nil
}
