package text

import (
	"fmt"
	"gokit/pkg/common"
	"gokit/pkg/i18n"
	conv2 "gokit/pkg/text/conv"
	count2 "gokit/pkg/text/count"
	"strings"
)

var (
	I18nKey_Description = i18n.I18nKey{Key: "text_description", Common: "文本相关工具"}
	description         = i18n.GetLocaleVal(I18nKey_Description)
)

type Text struct {
}

func (text *Text) Description() string {
	return description
}
func (text *Text) Run(args []string) {
	opType, args := common.VisitOne(args)
	switch strings.ToLower(opType) {
	case "count":
		count(args)
		break
	case "conv":
		conv(args)
		break
	default:
	}
}
func printUsage() {
	fmt.Println("文本工具，可进行统计转化等操作")
	fmt.Println("使用方式 操作[count,conv] 参数")
}
func count(args []string) {
	count2.Run(args)
}
func conv(args []string) {
	conv2.Run(args)
}
