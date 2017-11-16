package main

import (
	"../gCrawl"
	//"time"
)

func main() {
	craw := gCrawl.NewCrawl(20)
	craw.Work("http://www.hao123.com")
	//time.Sleep(time.Second * 5)
	//craw.Stop()
	select{}
}
