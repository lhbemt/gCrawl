package gParseLinks

import (
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

func parseDetail(mainurl string, keyword string,urls []string, results []string, node *html.Node)([]string, []string) {
	if node.Type == html.ElementNode && node.Data == "a" { // find releation links
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				var bFind = false
				bDirect := strings.HasPrefix(attr.Val, mainurl)
				if bDirect {
					urls = append(urls, attr.Val)
				} else {
					if strings.HasPrefix(attr.Val, keyword) {
						bFind = true
					}
					if attr.Val[0] == '/' {
						url := mainurl + attr.Val
						urls = append(urls, url)
					}
				}
				if bFind {
					var title []string
					title = append(title, node.Data, attr.Val)
					result := strings.Join(title, " ")
					results = append(results, result)
				}
			}
		}
	}

	return urls, results
}

func parseHtml(mainurl string, keyword string,urls []string, results []string, node *html.Node)([]string, []string) {
	urls, results = parseDetail(mainurl, keyword, urls, results, node)
	for ch := node.FirstChild; ch != nil; ch = ch.NextSibling {
		urls, results = parseHtml(mainurl, keyword, urls, results, ch)
	}
	return urls, results
}

func ParseLinks(url string, header string, keyword string)([]string,[]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	} else {
		node, err2 := html.Parse(resp.Body)
		defer resp.Body.Close()
		if err2 != nil {
			return nil, nil,  err2
		} else {
			urls, results := parseHtml(url, keyword,nil, nil, node)
			return urls, results, nil
		}
	}
}
