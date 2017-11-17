package main

import (
	"../gCrawl"
	//"time"
	"../gConfigFile"
	"fmt"
)

func main() {
	mainurl, err := gConfigFile.GetMainUrl()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	header, err := gConfigFile.GetHeader()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	keyWord, err := gConfigFile.GetKeyWord()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	routineNum, err := gConfigFile.GetRoutineNum()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	craw := gCrawl.NewCrawl(routineNum)
	craw.Work(mainurl, header, keyWord)
	//time.Sleep(time.Second * 5)
	//craw.Stop()
	select{}
}
