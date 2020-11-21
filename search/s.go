package search

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

var searchURL = "http://zimuku.la/search?q="

type List struct {
	Keyword string // 搜索关键字
	Content []Content
}

type Content struct {
	ID     string // ID供下载使用
	Title  string
	Link   string
	DLNums float64 // DLNums为下载次数,判断字幕质量
}

// Get根据List的Keyword返回一个包含字幕Content列表的List
func (l *List) Get() (*List, error) {

	res, err := http.Get(searchURL + l.Keyword)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	list := new(List)

	var parseErr error
	doc.Find("div>table>tbody>tr[class!=msub]").Each(func(i int, s *goquery.Selection) {
		srtTag := s.Find("a")
		link, _ := srtTag.Attr("href")
		id := regexp.MustCompile(`/detail/([0-9]*?)\.html`).FindSubmatch([]byte(link))
		title := srtTag.Text()
		dlTag := s.Find("td[class=last]>div[class=hidden-xs]").Text()
		var dlNums float64

		r := []rune(dlTag)
		if string(r[len(r)-1]) == "万" {
			n, err := strconv.ParseFloat(string(r[:len(r)-1]), 64)
			if err != nil {
				parseErr = err
				return
			}
			dlNums = n * 10000
		} else {
			dlNums, err = strconv.ParseFloat(string(r), 64)
			if err != nil {
				parseErr = err
				return
			}
		}

		c := &Content{}
		c.Title = title
		c.Link = link
		c.DLNums = dlNums
		c.ID = string(id[1])

		list.Content = append(list.Content, *c)

	})

	if parseErr != nil {
		return nil, parseErr
	}
	return list, nil

}

// GetMostDownloads 获取最多下载次数的字幕Content
func (l *List) GetMostDownloads() (*Content, error) {
	mostDownload := &Content{
		DLNums: 0,
	}
	c := l.Content

	for _, item := range c {
		if item.DLNums >= mostDownload.DLNums {
			*mostDownload = item
		}
	}

	return mostDownload, nil
}
