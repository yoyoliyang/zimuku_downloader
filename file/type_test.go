package file

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

var bs = make([]byte, 32)

func TestFileType(t *testing.T) {
	tests := [...]string{"rar", "zip", "ass"}
	for _, v := range tests {
		f := "test" + "." + v
		buf, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatal(err)
		}

		var result bool
		switch v {
		case "rar":
			result = IsRar(buf)
		case "ass":
			result = IsAss(buf)
		case "zip":
			result = IsZip(buf)
		}

		if result != true {
			t.Errorf("%v should be %v type file", f, v)
		} else {
			fmt.Printf("%v ok\n", f)
		}

	}
}
