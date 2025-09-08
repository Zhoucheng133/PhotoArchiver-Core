package main

import "C"

/*
#include <stdlib.h>
#include <string.h>
*/
import (
	"encoding/json"
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
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil || x == nil {
		return ""
	}
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

//export ScanDir
func ScanDir(path *C.char) *C.char {
	files := scanDir(C.GoString(path))
	if len(files) == 0 {
		return C.CString("[]")
	}
	data, err := json.Marshal(files)
	if err != nil {
		return C.CString("[]")
	}
	return C.CString(string(data))
}

func scanDir(path string) []Photo {
	var files []Photo
	entries, err := os.ReadDir(path)
	if err != nil {
		return files
	}
	for _, entry := range entries {
		if !entry.IsDir() && entry.Name() != ".DS_Store" {
			datetime := getCaptureTime(filepath.Join(path, entry.Name()))
			if datetime != "" {
				files = append(files, Photo{
					Name:     entry.Name(),
					Dir:      path,
					DateTime: datetime,
				})
			}
		}
	}

	return files
}

func main() {
	fmt.Println(scanDir("/Users/zhoucheng/Downloads/照片"))
}
