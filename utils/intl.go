package utils

import (
	"strings"
	"sync"

	"github.com/liuzl/gocc"
)

var (
	t2sConverter *gocc.OpenCC
	s2tConverter *gocc.OpenCC
	convertOnce  sync.Once
)

func initConverters() {
	convertOnce.Do(func() {
		t2sConverter, _ = gocc.New("t2s")
		s2tConverter, _ = gocc.New("s2t")
	})
}

func containsHan(s string) bool {
	for _, r := range s {
		if r >= 0x4e00 && r <= 0x9fff {
			return true
		}
	}
	return false
}

func ConvertCountry(name string) map[string]string {
	nameLower := strings.ToLower(name)
	return map[string]string{
		"enUS": name,
		"zhCN": countryNameCN[nameLower],
		"zhTW": countryNameTW[nameLower],
	}
}

func ConvertCity(nameAlt string, name string) map[string]string {

	initConverters()
	names := strings.Split(nameAlt, ",")

	var hanName, firstName string

	for _, raw := range names {
		n := strings.TrimSpace(raw)
		if n == "" {
			continue
		}
		if firstName == "" {
			firstName = n
		}
		if hanName == "" && containsHan(n) {
			hanName = n
		}
	}
	cnName, twName := hanName, hanName
	if hanName != "" {
		if t2sConverter != nil {
			if s, err := t2sConverter.Convert(hanName); err == nil {
				cnName = s
			}
		}
		if s2tConverter != nil {
			if s, err := s2tConverter.Convert(hanName); err == nil {
				twName = s
			}
		}
	} else {
		// 完全没有中文名，兜底用英文名占位
		cnName = name
		twName = name
	}

	return map[string]string{
		"enUS": name,
		"zhCN": cnName,
		"zhTW": twName,
	}
}
