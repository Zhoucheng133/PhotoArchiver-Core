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
	"photoarchiver/types"
	"photoarchiver/utils"
	"sync/atomic"
)

var stopFlag atomic.Bool

//export GetPhoto
func GetPhoto(path *C.char) *C.char {
	var data types.Photo = utils.GetPhoto(C.GoString(path))
	if data == (types.Photo{}) {
		return C.CString("")
	}

	json, _ := json.Marshal(data)

	return C.CString(string(json))
}

//export StopScan
func StopScan() {
	stopFlag.Store(true)
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

func scanDir(path string) []types.Photo {
	stopFlag.Store(false)
	var files []types.Photo
	entries, err := os.ReadDir(path)
	if err != nil {
		return files
	}
	for _, entry := range entries {
		if stopFlag.Load() {
			return files
		}
		if !entry.IsDir() && entry.Name() != ".DS_Store" {
			photoPath := filepath.Join(path, entry.Name())
			photo := utils.GetPhoto(photoPath)
			if photo.DateTime != "" {
				files = append(files, photo)
			}
		}
	}

	return files
}

func main() {
	fmt.Println(scanDir("/Users/zhoucheng/Downloads/照片"))
}
