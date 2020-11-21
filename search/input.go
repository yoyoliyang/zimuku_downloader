package search

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	//	"strings"
)

// Choose 输入字符串，返回0（回车）或数字（所选的）
func Choose(rang uint8) uint8 {
	s := bufio.NewScanner(os.Stdin)
	for {
		var input uint8
		for s.Scan() {
			// 处理回车\n符号
			// text := strings.TrimSuffix(s.Text(), "\n")
			// fmt.Print(text)
			text := s.Text()

			if text == "" {
				return 0
			}
			if i, err := strconv.ParseUint(text, 10, 32); err == nil {
				input = uint8(i)
				break
			}
		}
		if input <= rang {
			return input
		}
		fmt.Print("error input, retry: ")
	}
}
