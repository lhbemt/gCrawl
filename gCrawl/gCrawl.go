package gCrawl

import (
	"errors"
	"sync"
	"../gParseLinks"
	"fmt"
	"time"
	"os"
)

type urlsQueue struct {
	queue []string
	guard sync.Mutex
}

type GCrawl struct {
	nRoutineCount int // routine count
	header string // url header
	keyword string // keyword to crawl
	close chan interface{} // close signal
	running bool
	urls *urlsQueue // urlqueue
	seenurls map[string](bool) // having seen urls
	urlready chan byte // url is ready
	localRecord chan string // write to local
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
					//fmt.Println(url)
					this.urls.guard.Unlock()
					geturls, getResults, err := gParseLinks.ParseLinks(url, this.header, this.keyword)
					if err != nil { // log
						fmt.Printf("parse links: %s error: %s", url, err)
						continue
					} else {
						this.urls.guard.Lock()
						for _, link := range geturls {
							_, ok := this.seenurls[link]
							if !ok {
								//fmt.Println(link)
								this.urls.queue = append(this.urls.queue, link)
								this.urlready <- '1' // ready for parse
							}
						}
						this.urls.guard.Unlock()
						for _, result := range getResults {
							this.localRecord <- result
						}
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

func (this *GCrawl) Work(mainUrl string, header string, keyword string) error {
	if this.running {
		return errors.New("crawl is working")
	} else {
		this.running = true
		this.urls = new(urlsQueue)
		this.urls.queue = append(this.urls.queue, mainUrl)
		this.seenurls = make(map[string]bool)
		this.urlready = make(chan byte, 1000000)
		this.localRecord = make(chan string, 1000)
		this.header = header
		this.keyword = keyword
		for i := 0; i < this.nRoutineCount; i++ {
			go this.crawl()
		}
		this.urlready <- '1' // start work
		go this.localRecordTxt() // write to local
	}
	return nil
}

func getFile() (*os.File, error) {
	timeNow := time.Now()
	fileFolder := fmt.Sprintf("%d%d%d", timeNow.Year(), timeNow.Month(), timeNow.Day())
	// if dir not exist, create
	err := os.Mkdir(fileFolder, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return nil, fmt.Errorf("create folder error: %v", err) // create failed
	}
	// create file
	var i int
	for {

		fileName := fmt.Sprintf("%s/%s_%d.txt", fileFolder, fileFolder, i)
		f, err := os.OpenFile(fileName, os.O_APPEND | os.O_WRONLY, 0)
		if err != nil && os.IsNotExist(err) {
			f, err = os.Create(fileName)
			if err != nil {
				return nil, fmt.Errorf("create file error: %v", err)
			} else {
				return f, nil
			}
		} else if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("create file error: %v", err)
		} else {
			finfo, err := f.Stat()
			if err != nil {
				i++
				f.Close()
				continue // change other file to write
			} else {
				if finfo.Size() > 100 * 1024 * 1024 {
					i++
					f.Close()
					continue // change other file to write
				} else {
					return f, nil
				}
			}
		}
	}
}

func (this *GCrawl) localRecordTxt() {
	var f *os.File
	var err error

	f, err = getFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	for url := range this.localRecord {
		finfo, err := f.Stat()
		if err != nil {
			fmt.Printf("get file info error: %v", err)
			f.Close()
			return
		} else {
			if finfo.Size() > 100 * 1024 * 1024 {
				f.Close() // write to other file
				f, err = getFile()
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
		url += "\n"
		_, err = f.WriteString(url)
		if err != nil {
			fmt.Println(err)
		}
	}

	if f != nil {
		f.Close()
	}
}

func (this *GCrawl) Stop() {
	close(this.close)
	close(this.localRecord)
	this.running = false
}
