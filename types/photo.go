package types

type Photo struct {
	Dir      string `json:"dir"`
	Name     string `json:"name"`
	DateTime string `json:"datetime"`
	Country  string `json:"country"`
	City     string `json:"city"`
}
