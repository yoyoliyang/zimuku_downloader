package main

import (
	"fmt"
	"log"
	"srt/search"
)

func main() {

	var content search.Content
	content.Keyword = "变形金刚"
	list, err := content.Get()
	if err != nil {
		log.Fatal(err)
	}

	h, err := list.GetMostDownloads()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(h.DLNums, h.Title, h.Link)

}
