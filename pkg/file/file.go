package file

import (
	"fmt"
	"gokit/pkg/common"
	"gokit/pkg/file/internal"
	"gokit/pkg/i18n"
)

var (
	I18nKey_Description = i18n.I18nKey{Key: "file_description", Common: "文件相关工具"}
	description         = i18n.GetLocaleVal(I18nKey_Description)
)

type File struct {
}

func printUsage() {
	fmt.Println("文件工具,使用方式 操作 参数")
	fmt.Println("操作 [cat,rm,mk,cp,touch,mv,ls,count]")
}
func (file *File) Description() string {
	return description
}
func (file *File) Run(args []string) {
	op, args := common.VisitOne(args)
	switch op {
	case "cat":
		internal.Cat(args)
	case "rm":
		internal.Rm(args)
	case "mk":
		internal.Mk(args)
	case "cp":
		internal.Cp(args)
	case "touch":
		internal.Touch(args)
	case "mv":
		internal.Mv(args)
	case "ls":
		internal.Ls(args)
	case "count":
		internal.Count(args)
	default:
		printUsage()
	}
}
