package image

import (
	"gokit/pkg/i18n"
)

var (
	I18nKey_Description = i18n.I18nKey{Key: "image_description", Common: "图片相关工具"}
	description         = i18n.GetLocaleVal(I18nKey_Description)
)

type Image struct {
}

func (image *Image) Description() string {
	return description
}
func (image *Image) Run(args []string) {
}
