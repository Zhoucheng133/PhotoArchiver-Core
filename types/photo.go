package types

type Photo struct {
	Dir      string            `json:"dir"`
	Name     string            `json:"name"`
	DateTime string            `json:"datetime"`
	Country  map[string]string `json:"country"`
	City     map[string]string `json:"city"`
}
