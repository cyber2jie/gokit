package i18n

import (
	"github.com/Xuanwo/go-locale"
	"github.com/go-ini/ini"
	"os"
	"path/filepath"
)

const (
	localePath       = "locales"
	localeFileSuffix = ".ini"
)

var i18nTable = map[string]string{}

type I18nKey struct {
	Key    string
	Common string
}

func init() {
	wd, error := os.Getwd()
	if error == nil {
		localeFilePath := wd + "/" + localePath + "/" + GetLocaleStr() + localeFileSuffix
		load(filepath.FromSlash(localeFilePath))
	}
}

//不管理section
func load(path string) {
	file, error := ini.Load(path)
	if error == nil {
		for _, section := range file.Sections() {
			for _, key := range section.Keys() {
				name := key.Name()
				val := key.Value()
				if name != "" && val != "" {
					i18nTable[name] = val
				}
			}
		}
	}
}
func GetLocaleStr() string {
	tag, _ := locale.Detect()
	return tag.String()
}
func GetLocaleVal(key I18nKey) string {
	localeVal := i18nTable[key.Key]
	if localeVal != "" {
		return localeVal
	}
	return key.Common
}
