package file

import (
	"fmt"
	"log"
	"os"

	"github.com/mholt/archiver/v3"
)

func ExtractFile(file string) error {
	defer fmt.Printf("%v > ./Subs\n", file)
	// archiver需要初始化结构体才能使用单独的配置（比如文件错误仍然解压或覆盖解压），但是此处需要处理zip和rar，单独创建两个结构体比较繁琐，故采用目录检测方式来判断是否已经存在了Subs目录（rargb下载的电影里面有Subs目录）
	_, err := os.Open("Subs")
	if err != nil {
		if os.IsNotExist(err) {
			err = archiver.Unarchive(file, "Subs")
			if err != nil {
				log.Fatalln(err)
			}
			return nil
		} else {
			log.Fatalln(err)
		}
	}

	err = os.RemoveAll("Subs")
	if err != nil {
		log.Fatalln(err)
	}
	err = archiver.Unarchive(file, "Subs")
	if err != nil {
		log.Fatalln(err)
	}

	return nil
}
