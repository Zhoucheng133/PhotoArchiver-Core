package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
)

func getCaptureTime(path string) string {
	f, _ := os.Open(path)
	defer f.Close()

	x, _ := exif.Decode(f)
	getTagString := func(name exif.FieldName) string {
		tag, err := x.Get(name)
		if err != nil || tag == nil {
			return ""
		}
		return tag.String()
	}
	originalTime := strings.ReplaceAll(getTagString(exif.DateTimeOriginal), "\"", "")
	return strings.Replace(originalTime, ":", "/", 2)
}

func archive() {

}

func main() {
	fmt.Println(getCaptureTime("/Users/zhoucheng/Downloads/DSC08068.ARW"))
}
