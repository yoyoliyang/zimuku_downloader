package main

import (
	"fmt"
	"log"
	"os"
	"srt/download"
	"srt/file"
	"srt/search"
)

// 流程：探测当前目录有无视频文件-获取文件名or手动输入-取得字幕列表-手动or自动选择-获取到编号-下载保存-解压or不解压-重命名
func main() {

	// 初始化搜索关键字列表
	list := &search.List{}

	videoName := &file.VideoName{}
	videoName, err := videoName.GetVideoInfo()
	if err != nil {
		log.Fatal(err)
	}

	if videoName.Name == "" {
		fmt.Print("当前目录下没有找到视频文件，手动输入: ")
		var input string
		if _, err := fmt.Scan(&input); err != nil {
			fmt.Println(err)
			return
		}
		list.Keyword = input
	} else {
		list.Keyword = videoName.Keyword
	}
	fmt.Println(list.Keyword)

	list, err = list.Get()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ID", "\t", "下载次数", "\t", "名称")
	for index, item := range list.Content {
		fmt.Println(index+1, "\t", item.DLNums, "\t", item.Title)
	}

	fmt.Print("输入编号进行下载or回车自动下载: ")
	// 调用Choose来处理输入
	cid := search.Choose(uint8(len(list.Content)))
	content := &search.Content{}
	if cid == 0 {
		content, err = list.GetMostDownloads()
	} else {
		for index, item := range list.Content {
			if cid == uint8(index+1) {
				*content = item
				break
			}
		}
	}
	fmt.Println(content)

	//  根据content处理下载
	dlContent := &download.DlContent{}
	dlContent.ID = content.ID
	SubData, err := dlContent.GetFile()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("downloading")

	// savedFile, err := os.Create(videoName.NameWithoutExt)
	savedFile, err := os.OpenFile(videoName.NameWithoutExt+content.Ext, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}

	_, err = savedFile.Write(SubData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("done.")

}
