package utils

import (
	"os"
	"path/filepath"
	"strings"

	"photoarchiver/types"

	"github.com/rwcarlsen/goexif/exif"
)

func init() {
	InitLocation()
}

func GetPhoto(path string) types.Photo {
	datetime := GetCaptureTime(path)
	if datetime == "" {
		return types.Photo{}
	}

	country, city := "", ""
	f, err := os.Open(path)
	if err == nil {
		defer f.Close()
		if x, err := exif.Decode(f); err == nil && x != nil {
			if lat, lng, err := x.LatLong(); err == nil {
				country, city = GetLocation(lat, lng)
			}
		}
	}

	return types.Photo{
		Dir:      filepath.Dir(path),
		Name:     filepath.Base(path),
		DateTime: datetime,
		Country:  country,
		City:     city,
	}
}

func GetCaptureTime(path string) string {
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
