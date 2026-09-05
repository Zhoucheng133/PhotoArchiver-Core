package utils

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"

	"photoarchiver/types"

	"errors"

	"github.com/andreiashu/geobed"
	"github.com/rwcarlsen/goexif/exif"
	"go4.org/media/heif"
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

func GetPhoto(path string) (types.Photo, error) {
	f, err := os.Open(path)
	if err != nil {
		return types.Photo{}, errors.New("cannot open")
	}
	defer f.Close()

	x, err := decodeEXIF(f, path)
	if err != nil || x == nil {
		return types.Photo{}, errors.New("no exif data")
	}

	datetime := getCaptureTime(x)
	country, city, cityAlt := getLocation(x)

	return types.Photo{
		Dir:      filepath.Dir(path),
		Name:     filepath.Base(path),
		DateTime: datetime,
		Country:  ConvertCountry(country),
		City:     ConvertCity(cityAlt, city),
	}, nil
}

func decodeEXIF(f *os.File, path string) (*exif.Exif, error) {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".heic", ".heif":
		return decodeHEIFEXIF(f)
	default:
		return exif.Decode(f)
	}
}

func decodeHEIFEXIF(f *os.File) (*exif.Exif, error) {
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	reader := io.NewSectionReader(f, 0, stat.Size())

	h := heif.Open(reader)

	rawEXIF, err := h.EXIF()
	if err != nil {
		return nil, err
	}

	return exif.Decode(bytes.NewReader(rawEXIF))
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

func getLocation(exifData *exif.Exif) (string, string, string) {
	country, city, cityAlt := "", "", ""
	lat, lng, err := exifData.LatLong()
	if err != nil {
		return country, city, cityAlt
	}
	if g == nil {
		InitLocation()
	}
	if g == nil {
		return country, city, cityAlt
	}
	result := g.ReverseGeocode(lat, lng)
	countryCode := result.Country()
	countryName := countryCode
	if name, ok := countryByISO[countryCode]; ok {
		countryName = name
		country = countryName
		city = result.City
		cityAlt = result.CityAlt
	}
	return country, city, cityAlt
}
