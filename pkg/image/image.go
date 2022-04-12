package image

import (
	"fmt"
	"gokit/pkg/common"
	"gokit/pkg/i18n"
)

var (
	I18nKey_Description = i18n.I18nKey{Key: "image_description", Common: "图片相关工具"}
	description         = i18n.GetLocaleVal(I18nKey_Description)
)
var actions = map[string]func([]string){
	"blur":      blur,
	"resize":    resise,
	"watermark": watermark,
}

type Image struct {
}

func printUsage() {
	fmt.Println("image action args")
	fmt.Println("action[blur,resise,watermark]")
}
func (image *Image) Description() string {
	return description
}
func (image *Image) Run(args []string) {
	action, args := common.VisitOne(args)
	actionFunc := actions[action]
	if actionFunc != nil {
		actionFunc(args)
	} else {
		printUsage()
	}
}
