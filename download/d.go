package download

import (
	"bytes"
	"errors"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

var dlURL = "http://zmk.pw"

type DlContent struct {
	ID        string
	DlURL     string
	StaticURL string
}

// GetFile获取到下载链接，并返回字节文件
func (d *DlContent) GetFile() ([]byte, error) {
	d.DlURL = dlURL + "/dld/" + d.ID + ".html"
	resp, err := http.Get(d.DlURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// 获取下载线路,如果匹配下面的线路，那么给staticURL赋值，如果匹配不到，为空
	doc.Find("div>table>tbody>tr>td>div>ul>li").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Find("a").Attr("href")
		// 联通线路下载
		dl := regexp.MustCompile(`.*/lt1`).Find([]byte(link))
		if dl != nil {
			d.StaticURL = dlURL + string(dl)
		}
	})

	// 没有匹配到时，出现了错误，返回错误信息
	if d.StaticURL == "" {
		return nil, errors.New("error for getting download static url")
	}

	req, err := http.NewRequest("GET", d.StaticURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Referer", d.DlURL)
	req.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36`)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data := &bytes.Buffer{}
	err = CopyWithProgress(data, resp)
	if err != nil {
		return nil, err
	}

	return data.Bytes(), nil

}
