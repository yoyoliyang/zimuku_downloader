package search

import (
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

var searchURL = "http://zimuku.la/search?q="

type List struct {
	Content []Content
}

type Content struct {
	Keyword string
	Title   string
	Link    string
	DLNums  float64
}

func (c *Content) Get() (*List, error) {

	res, err := http.Get(searchURL + c.Keyword)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	list := new(List)

	doc.Find("div>table>tbody>tr[class!=msub]").Each(func(i int, s *goquery.Selection) {
		srtTag := s.Find("a")
		link, _ := srtTag.Attr("href")
		title := srtTag.Text()
		dlTag := s.Find("td[class=last]>div[class=hidden-xs]").Text()
		var dlNums float64

		r := []rune(dlTag)
		if string(r[len(r)-1]) == "万" {
			n, _ := strconv.ParseFloat(string(r[:len(r)-1]), 64)
			dlNums = n * 10000
		} else {
			dlNums, _ = strconv.ParseFloat(string(r), 64)
		}

		c.Title = title
		c.Link = link
		c.DLNums = dlNums

		list.Content = append(list.Content, *c)

	})

	return list, nil

}

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
