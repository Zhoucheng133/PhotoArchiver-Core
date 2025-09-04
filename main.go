package main

import "C"

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
)

type Photo struct {
	Dir      string `json:"dir"`
	Name     string `json:"name"`
	DateTime string `json:"datetime"`
}

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

func scanDir(path string) []Photo {
	var files []Photo
	entries, err := os.ReadDir(path)
	if err != nil {
		return files
	}
	for _, entry := range entries {
		if !entry.IsDir() && entry.Name() != ".DS_Store" {
			files = append(files, Photo{
				Name:     entry.Name(),
				Dir:      path,
				DateTime: getCaptureTime(filepath.Join(path, entry.Name())),
			})
		}
	}

	return files
}

func main() {
	fmt.Println(scanDir("/Users/zhoucheng/Pictures/临时"))
}
