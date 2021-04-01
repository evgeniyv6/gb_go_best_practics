package main

import "log"

func goodHash(done <-chan interface{}, files <-chan string) <-chan interface{} {
	completed := make(chan interface{})

	go func() {
		defer log.Println("terminated child")
		defer close(completed)
		for {
			select {
			case <-done:
				return
			case f := <-files:
				h, err := compute(f)
				if err != nil {
					log.Printf("compute hash for %q: %v", f, err)
					continue
				}
				log.Printf("hash for file %q: %q", f, h)

			}
		}

	}()

	return completed
}
