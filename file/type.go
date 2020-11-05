package file

import (
	"bytes"
	"strings"
)

// GetFileType 获取文件的类型，判断方式如下：
// ass 文件类型，取第一行数据，如果包含Script Info关键字，那么为ass格式
// 因为下载的大部分ass文件不是utf8编码(golang无法匹配不是统一编码的字符串或字节)，那么对前面的几个字节逐字分析，匹配成[Script(ass文件内容)的话，为ass格式
// rar及zip文件取file magic number

var (
	zipHeader = [...]byte{'\x50', '\x4b', '\x03', '\x04'}
	rarHeader = [...]byte{'\x52', '\x61', '\x72', '\x21', '\x1A', '\x07', '\x00'}
	ass       = "[Script"
)

func IsAss(bs []byte) bool {
	res := make([]byte, 1)
	for _, v := range string(bs[:16]) {
		switch string(v) {
		case "[":
			res = append(res[:], byte(v))
		case "S":
			res = append(res[:], byte(v))
		case "c":
			res = append(res[:], byte(v))
		}
	}

	if strings.HasPrefix(ass, string(res)) {
		return true
	}
	return false
}

func IsRar(bs []byte) bool {
	if bytes.HasPrefix(bs[:32], rarHeader[:]) {
		return true
	}
	return false
}

func IsZip(bs []byte) bool {
	if bytes.HasPrefix(bs[:32], zipHeader[:]) {
		return true
	}
	return false
}
