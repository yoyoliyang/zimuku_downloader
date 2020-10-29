package file

import (
	"fmt"
	"os"
	"regexp"
)

type VideoName struct {
	Name           string
	Keyword        string
	NameWithoutExt string
}

// GetVideoName返回当前目录的视频文件名称
// 如果当前目录不包含视频文件，不进行错误处理，交由主函数进行手动处理，下同
func (v *VideoName) GetVideoInfo() (*VideoName, error) {
	// 此处需要初始化v，避免其他函数调用时，无法处理错误
	v = &VideoName{
		Name:           "",
		Keyword:        "",
		NameWithoutExt: "",
	}
	f, err := os.Open(".")
	if err != nil {
		return v, err
	}
	fList, err := f.Readdir(0)
	if err != nil {
		return v, err
	}

	for _, item := range fList {
		n := item.Name()
		if len(n) > 4 {
			if n[len(n)-4:] == ".mp4" || n[len(n)-4:] == ".mkv" {
				v.Name = n
				v.NameWithoutExt = n[:len(n)-4]
			}
		}
	}

	// 搜索视频文件，并且赋值v视频文件的搜索关键字，根据rarbg视频文件命名，取视频发行年代前的字符
	re := regexp.MustCompile(`(.*?)\.(\d{4}?)\..*\.(mp4|mkv)`)
	result := re.FindStringSubmatch(v.Name)
	if len(result) == 0 {
		return v, nil
	}

	fmt.Println("发现视频文件：", result[0])
	fmt.Println("搜索关键字：", result[1])
	v.Keyword = result[1]
	return v, nil

}
