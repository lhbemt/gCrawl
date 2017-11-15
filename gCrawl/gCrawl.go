package gCrawl

import (
	"errors"
	"sync"
	"../gParseLinks"
	"fmt"
)

type urlsQueue struct {
	queue []string
	guard sync.Mutex
}

type GCrawl struct {
	nRoutineCount int // routine count
	close chan interface{} // close signal
	running bool
	urls *urlsQueue // urlqueue
	seenurls map[string](bool) // having seen urls
	urlready chan byte // url is ready
}

func NewCrawl(nRoutineCount int) (*GCrawl) {
	if nRoutineCount <= 0 {
		nRoutineCount = 20
	}
	return &GCrawl{nRoutineCount: nRoutineCount, close: make(chan interface{}), running: false}
}

func (this *GCrawl) crawl() {
	for {
		select {
		case <-this.close : // terminate routine
			return
		case <-this.urlready : // read url and find links
			this.urls.guard.Lock()
			if len(this.urls.queue) > 0 {
				url := this.urls.queue[0]
				_, ok := this.seenurls[url]
				if !ok {
					this.seenurls[url] = true
					this.urls.queue = this.urls.queue[1:]
					this.urls.guard.Unlock()
					geturls, err := gParseLinks.ParseLinks(url)
					if err != nil { // log
						fmt.Printf("parse links: %s error: %s", url, err)
						continue
					} else {
						this.urls.guard.Lock()
						for _, link := range geturls {
							_, ok := this.seenurls[link]
							if !ok {
								fmt.Println(link)
								this.urls.queue = append(this.urls.queue, link)
								this.urlready <- '1' // ready for parse
							}
						}
						this.urls.guard.Unlock()
						continue
					}
				} else {
					this.urls.queue = this.urls.queue[1:]
					this.urls.guard.Unlock()
					continue
				}
			}
			this.urls.guard.Unlock()
		}
	}
}

func (this *GCrawl) Work(mainUrl string) error {
	if this.running {
		return errors.New("crawl is working")
	} else {
		this.running = true
		this.urls = new(urlsQueue)
		this.urls.queue = append(this.urls.queue, mainUrl)
		this.seenurls = make(map[string]bool)
		this.urlready = make(chan byte, 1000000)
		for i := 0; i < this.nRoutineCount; i++ {
			go this.crawl()
		}
		this.urlready <- '1' // start work
	}
	return nil
}

func (this *GCrawl) Stop() {
	close(this.close)
	this.running = false
}
