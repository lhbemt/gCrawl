package gParseLinks

import (
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

func parseHtml(mainurl string, urls []string, node *html.Node)([]string) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				bDirect := strings.Contains(attr.Val, "http")
				if bDirect {
					urls = append(urls, attr.Val)
				} else {
					link := mainurl + "/" + attr.Val
					urls = append(urls, link)
				}
			}
		}
	}
	for ch := node.FirstChild; ch != nil; ch = ch.NextSibling {
		urls = parseHtml(mainurl, urls, ch)
	}
	return urls
}

func ParseLinks(url string)([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	} else {
		node, err2 := html.Parse(resp.Body)
		if err2 != nil {
			return nil, err2
		} else {
			return parseHtml(url,nil, node), nil
		}
	}
}
