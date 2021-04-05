package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type crawlRes struct {
	err error
	msg string
}

type crawler struct {
	sync.Mutex
	visited map[string]string
	maxDepth int
}

func newCrawler(maxDepth int) *crawler {
	return &crawler{
		visited: make(map[string]string),
		maxDepth: maxDepth,
	}
}

func (c *crawler) run(ctx context.Context, url string, res chan <- crawlRes, depth int) {
	time.Sleep(2*time.Second)
	ctxDeadline, cancelFunc := context.WithDeadline(context.Background(), time.Now().Add(timeout))
	defer cancelFunc()

	select {
	case <-ctx.Done():
		return
	default:
		c.Lock()
		if depth >= c.maxDepth {
			return
		}
		c.Unlock()

		page, err := parse(ctxDeadline, url)

		if err != nil {
			res <- crawlRes{
				err: errors.New(fmt.Sprintf("err: %v when parse - %s", err, page)),
			}
			return
		}

		title := pageTitle(ctxDeadline, page)
		links := pageLinks(ctxDeadline, nil, page)

		c.Lock()
		c.visited[url] = title
		c.Unlock()

		res <- crawlRes{
			err: nil,
			msg: fmt.Sprintf("%s -> %s\n", url, title),
		}

		for l := range links {
			if c.checkVisited(l) {
				continue
			}
			c.Lock()
			depth++  // !! иначе depth >= c.maxDepth всегда true
			c.Unlock()
			go c.run(ctx,l, res,depth)
		}

	}
}

func (c *crawler) checkVisited(url string) bool {
	c.Lock()
	defer c.Unlock()

	_, ok := c.visited[url]
	return ok
}