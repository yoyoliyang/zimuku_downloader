package file

import (
	"fmt"
	"os"
	"regexp"

	"github.com/mholt/archiver/v3"
	"github.com/pkg/errors"
)

var err = errors.New("Err >>> file/archive.go:")

func ExtractFile(file string) error {
	err := archiver.Unarchive(file, "Subs")
	if err != nil {
		return errors.Wrap(err, err.Error())
	}

	f, err := os.Open("Subs")
	if err != nil {
		return err
	}
	info, err := f.Readdir(1)
	if err != nil {
		return err
	}

	for _, v := range info {
		if r, err := regexp.MatchString(`.*[chs|CHS|简体].*`, v.Name()); r == true {
			if err != nil {
				return errors.Wrap(err, err.Error())
			}
			fmt.Println(v.Name())
		}
	}

	return nil
}
