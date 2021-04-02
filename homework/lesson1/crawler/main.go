package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	errLimit = 100000
	resLimit = 10000
	timeout time.Duration = 1 * time.Second
)

var (
	url string
	depthLimit int
)


func init() {
	flag.StringVar(&url,"url","", "url address")
	flag.IntVar(&depthLimit, "depth", 3, "max depth for run")
	flag.Parse()

	if url == "" {
		log.Println("url address is empty as parameter")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	started := time.Now()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	crawler := newCrawler(depthLimit)
	go watchSignals(cancel, crawler)

	res := make(chan crawlRes)
	done := watchCrawler(ctx, res, errLimit, resLimit)
	crawler.run(ctx, url,res,0)

	<-done
	log.Println(time.Since(started))

}

func watchSignals(cancel context.CancelFunc, cr *crawler) {
	osSignalChan := make(chan os.Signal)
	signal.Notify(osSignalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	for {
		sig := <-osSignalChan
		switch sig {
		case syscall.SIGUSR1:
			log.Println("increased by 10")
			cr.maxDepth += 10
		case syscall.SIGINT, syscall.SIGTERM:
			log.Printf("got signal %q", sig.String())
			cancel()
		}
	}
}

func watchCrawler(ctx context.Context, res <- chan crawlRes, maxErr, maxRes int) chan struct{} {
	readersDone := make(chan struct{})

	go func() {
		defer close(readersDone)
		for {
			select {
			case <- ctx.Done():
				return
			case r := <- res:
				if r.err != nil {
					maxErr--
					if maxErr <= 0 {
						log.Println("max errors exceeded")
						return
					}
					continue
				}
				log.Printf("crawling result: %v", r.msg)
				maxRes--
				if maxRes <= 0 {
					log.Println("got max results")
					return
				}
			}
		}
	}()
	return readersDone
}
