package download

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// CopyWithProgress 拷贝http相应的Body到buf中，并显示进度
func CopyWithProgress(dst *bytes.Buffer, src *http.Response) error {

	count := make(chan int)
	done := make(chan bool)
	err := make(chan error)

	// 创建一个缓冲区，写入数据，并记录数据量
	go func() {
		for {
			count <- dst.Len()
			time.Sleep(time.Nanosecond)
		}

	}()

	go func() {
		for {
			select {
			case c := <-count:
				fmt.Print("\r", c, "/", src.ContentLength)
				// 当计数器和写入的数据大小相同时，给done通道赋值
				if int64(c) == src.ContentLength {
					fmt.Print("\r", c, "/", src.ContentLength, "\n")
					done <- true
				}
				time.Sleep(time.Second / 2)
			}
		}
	}()

	go func() {
		_, e := io.Copy(dst, src.Body)
		if e != nil {
			err <- e
			return
		}
		err <- nil
	}()

	// data.Reset()
	<-done
	fmt.Println("done")

	return <-err
}
