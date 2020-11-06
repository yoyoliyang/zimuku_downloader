package file

import (
	"log"
	"testing"
)

func TestExtractFile(t *testing.T) {
	f := "test.rar"

	err := ExtractFile(f)
	if err != nil {
		log.Fatal(err)
	}
}
