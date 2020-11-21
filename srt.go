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
	searchList := &search.List{}

	var videoName = &file.VideoName{}

	// 手动搜索字幕时的方法,将videoName直接重新赋值，不管当前有没有视频文件
	args := os.Args
	if len(args) == 2 {
		if args[1] == "-h" {
			fmt.Printf("zimuku.la下载工具:\n\t输入 %v 视频名称，或只输入%v自动探测当前目录视频文件\n", args[0], args[0])
			return
		}
		videoName = &file.VideoName{
			Name:           args[1],
			Keyword:        args[1],
			NameWithoutExt: args[1],
		}
	} else {
		// 自动探测当前目录视频文件
		vn, err := videoName.GetVideoInfo()
		if err != nil {
			log.Fatal(err)
		}
		if vn.Name != "" {
			videoName = vn
		}
	}

	// 找不到视频文件时的处理
	if videoName.Name == "" {
		fmt.Print("当前目录下没有找到视频文件，手动输入: ")
		var input string
		if _, err := fmt.Scan(&input); err != nil {
			fmt.Println(err)
			return
		}
		videoName = &file.VideoName{
			Name:           input,
			Keyword:        input,
			NameWithoutExt: input,
		}
	}

	searchList.Keyword = videoName.Keyword
	fmt.Println(searchList.Keyword)

	searchList, err := searchList.Get()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ID", "\t", "下载次数", "\t", "名称")
	for index, item := range searchList.Content {
		fmt.Println(index+1, "\t", item.DLNums, "\t", item.Title)
	}

	fmt.Print("输入编号进行下载or回车自动下载: ")
	// 调用Choose来处理输入
	cid := search.Choose(uint8(len(searchList.Content)))
	content := &search.Content{}
	if cid == 0 {
		content, err = searchList.GetMostDownloads()
	} else {
		for index, item := range searchList.Content {
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
	bs, err := dlContent.GetFile()
	if err != nil {
		log.Fatal(err)
	}

	// 获取到最终下载的文件字节切片后，进行文件类型检测，并根据文件类型来具体操作（解压或重命名）
	var savedFileName = videoName.NameWithoutExt
	switch {
	case file.IsAss(bs) == true:
		savedFileName += ".ass"
	case file.IsRar(bs) == true:
		savedFileName += ".rar"
		defer file.ExtractFile(savedFileName)
	case file.IsZip(bs) == true:
		savedFileName += ".zip"
		defer file.ExtractFile(savedFileName)
	default:
		savedFileName += ".srt"
	}

	savedFile, err := os.OpenFile(savedFileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}

	_, err = savedFile.Write(bs)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(content.Title, ">", savedFileName)

}
