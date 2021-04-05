package main

import (
	"context"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
)

func parse(ctx context.Context, url string) (*html.Node, error) {
	select {
	case <-ctx.Done():
		log.Println(ctx.Err())
		return nil,ctx.Err()
	default:
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("cannot get page")
		}
		b, err := html.Parse(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("can't parse page")
		}
		return b, err
	}
}

func pageTitle(ctx context.Context, n *html.Node) string {
	var title string

	select {
	case <-ctx.Done():
		log.Println(ctx.Err())
	default:
		if n.Type == html.ElementNode && n.Data == "title" {
			return n.FirstChild.Data
		}

		for c:=n.FirstChild; c != nil; c = c.NextSibling {
			title = pageTitle(ctx, c)
			if title != "" {
				break
			}
		}
	}

	return title
}

func pageLinks(ctx context.Context, links map[string]struct{}, n *html.Node) map[string]struct{} {
	if links == nil {
		links = make(map[string]struct{})
	}

	select {
	case <-ctx.Done():
		log.Println(ctx.Err())
	default:
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}

				if _, ok := links[a.Val]; !ok && len(a.Val) >2 && a.Val[:2]=="//" {
					links["http://"+a.Val[2:]] = struct{}{}
				}
			}
		}

		for c := n.FirstChild; c!= nil ; c =c.NextSibling {
			links = pageLinks(ctx, links,c)
		}
	}

	return links
}
