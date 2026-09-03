package utils

import (
	"github.com/andreiashu/geobed"
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

// lat为纬度, lng为经度
func GetLocation(lat float64, lng float64) (string, string) {
	result := g.ReverseGeocode(lat, lng)
	countryCode := result.Country()
	countryName := countryCode
	if name, ok := countryByISO[countryCode]; ok {
		countryName = name
	}
	return countryName, result.City
}
