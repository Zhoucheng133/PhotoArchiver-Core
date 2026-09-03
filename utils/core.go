package utils

import (
	"os"
	"path/filepath"
	"strings"

	"photoarchiver/types"

	"github.com/andreiashu/geobed"
	"github.com/rwcarlsen/goexif/exif"
)

var g *geobed.GeoBed
var countryByISO map[string]string

func InitLocation() {
	var err error
	g, err = geobed.NewGeobed()
	if err != nil {
		return
	}
	countryByISO = make(map[string]string, len(g.Countries))
	for _, c := range g.Countries {
		countryByISO[c.ISO] = c.Country
	}
}

func GetPhoto(path string) types.Photo {
	f, err := os.Open(path)
	if err != nil {
		return types.Photo{}
	}
	defer f.Close()
	x, err := exif.Decode(f)
	if err != nil || x == nil {
		return types.Photo{}
	}

	datetime := getCaptureTime(x)
	country, city := getLocation(x)

	return types.Photo{
		Dir:      filepath.Dir(path),
		Name:     filepath.Base(path),
		DateTime: datetime,
		Country:  country,
		City:     city,
	}
}

func getCaptureTime(exifData *exif.Exif) string {
	getTagString := func(name exif.FieldName) string {
		tag, err := exifData.Get(name)
		if err != nil || tag == nil {
			return ""
		}
		return tag.String()
	}
	originalTime := strings.ReplaceAll(getTagString(exif.DateTimeOriginal), "\"", "")
	return strings.Replace(originalTime, ":", "/", 2)
}

func getLocation(exifData *exif.Exif) (string, string) {
	country, city := "", ""
	lat, lng, err := exifData.LatLong()
	if err != nil {
		return country, city
	}
	if g == nil {
		InitLocation()
	}
	if g == nil {
		return country, city
	}
	result := g.ReverseGeocode(lat, lng)
	countryCode := result.Country()
	countryName := countryCode
	if name, ok := countryByISO[countryCode]; ok {
		countryName = name
		country = countryName
		city = result.City
	}
	return country, city
}
