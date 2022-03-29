package internal

import (
	"fmt"
	"gokit/pkg/common"
	"io/ioutil"
	"os"
	"strings"
)

const (
	type_FILE = iota
	type_DIR
)

type typeInfo struct {
	tp   int
	name string
}

func writeError(builder strings.Builder, info string, err error) {
	builder.WriteString("[")
	builder.WriteString(info)
	builder.WriteString("]")
	builder.WriteString(err.Error())
	builder.WriteString("\r\n")
}
func isEmpty(array []string) bool {
	return len(array) == 0
}
func Cat(args []string) {
	if isEmpty(args) {
		fmt.Println("cat file")
		return
	}
	file, args := common.VisitOne(args)
	content, error := ioutil.ReadFile(file)
	if error != nil {
		fmt.Println("[", file, "]", error.Error())
		return
	}
	print := strings.Builder{}
	print.Write(content)
	fmt.Println(print.String())
}
func Rm(args []string) {
	if isEmpty(args) {
		fmt.Println("rm file1 file2 ....")
		return
	}
	err := strings.Builder{}
	for _, path := range args {
		error := os.Remove(path)
		if error != nil {
			writeError(err, path, error)
		}
	}
	fmt.Println(err.String())
}
func Mk(args []string) {
	if isEmpty(args) {
		fmt.Println("mk dir1 dir2 ....")
		return
	}
	err := strings.Builder{}
	for _, path := range args {
		error := os.MkdirAll(path, 0777)
		if error != nil {
			writeError(err, path, error)
		}
	}
	fmt.Println(err.String())
}
func Touch(args []string) {
	if isEmpty(args) {
		fmt.Println("touch file1 file2 ....")
		return
	}
	err := strings.Builder{}
	for _, file := range args {
		touched, error := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0666)
		if error != nil {
			writeError(err, file, error)
		} else {
			touched.Close()
		}
	}
	fmt.Println(err.String())
}
func Mv(args []string) {
	if isEmpty(args) {
		fmt.Println("mv oldfile newfile")
		return
	}
	oldFile, args := common.VisitOne(args)
	newFile, args := common.VisitOne(args)
	if newFile == "" {
		fmt.Println("mv need newFile argument")
		return
	}
	error := os.Rename(oldFile, newFile)
	if error != nil {
		fmt.Println(error.Error())
		return
	}
}
func Cp(args []string) {
	if isEmpty(args) {
		fmt.Println("cp oldfile copyfile")
		return
	}
	oldFile, args := common.VisitOne(args)
	newFile, args := common.VisitOne(args)
	if newFile == "" {
		fmt.Println("cp need copyFile argument")
		return
	}
	error := os.Link(oldFile, newFile)
	if error != nil {
		fmt.Println(error.Error())
		return
	}
}
func Ls(args []string) {
	if isEmpty(args) {
		fmt.Println("ls dir")
		return
	}
	dirPath, args := common.VisitOne(args)
	dirInfo, error := os.Stat(dirPath)
	if error != nil {
		fmt.Println("[", dirPath, "]", error.Error())
		return
	}
	if !dirInfo.IsDir() {
		fmt.Println(dirPath)
	} else {
		dir, error := os.Open(dirPath)
		if error != nil {
			fmt.Println("[", dirPath, "]", error.Error())
			return
		}
		dirNames, error := dir.Readdirnames(0)
		if error != nil {
			fmt.Println("[", dirPath, "]", error.Error())
			return
		}
		for _, dirName := range dirNames {
			fmt.Println(dirName)
		}
	}
}
func Count(args []string) {
	if isEmpty(args) {
		fmt.Println("count dir options")
		fmt.Println("options [-s suffix1 suffix2....] ")
		return
	}
	dirPath, args := common.VisitOne(args)
	option, args := common.VisitOne(args)
	dirInfo, error := os.Stat(dirPath)
	if error != nil {
		fmt.Println(error.Error())
		return
	}
	if !dirInfo.IsDir() {
		fmt.Println(dirPath, "is not a directory")
		return
	}
	suffixArray := []string{}
	suffixCountMap := map[string]int{}
	switch option {
	case "-s":
		for _, suffix := range args {
			contain := false
			for _, suffixInArray := range suffixArray {
				if suffix == suffixInArray {
					contain = true
					break
				}
			}
			if !contain {
				suffixArray = append(suffixArray, suffix)
				suffixCountMap[suffix] = 0
			}
		}
	default:
	}
	dir, error := os.Open(dirPath)
	if error != nil {
		fmt.Println(error.Error())
		return
	}
	typeinfos := list(dir, dirPath)
	fileCount, dirCount := 0, 0
	for _, info := range typeinfos {
		switch info.tp {
		case type_DIR:
			dirCount++
		case type_FILE:
			fileCount++
			if len(suffixArray) > 0 {
				for _, suffix := range suffixArray {
					{
						if common.EndWith(info.name, suffix) {
							suffixCount := suffixCountMap[suffix]
							suffixCount++
							suffixCountMap[suffix] = suffixCount
						}
					}
				}
			}
		}
	}
	fmt.Println("total ", dirCount, " dir")
	fmt.Println("total ", fileCount, " file")
	for suffix, count := range suffixCountMap {
		fmt.Println("total ", count, " ", suffix)
	}
}
func list(file *os.File, prefix string) []typeInfo {
	typeInfos := []typeInfo{}
	fileInfos, error := file.Readdir(0)
	if error == nil {
		for _, fileInfo := range fileInfos {
			typeInfo := typeInfo{name: fileInfo.Name()}
			if fileInfo.IsDir() {
				typeInfo.tp = type_DIR
				dirPath := prefix + "/" + fileInfo.Name()
				dir, error := os.Open(dirPath)
				if error == nil {
					subTypeInfos := list(dir, dirPath)
					typeInfos = append(typeInfos, subTypeInfos...)
				}
			} else {
				typeInfo.tp = type_FILE
			}
			typeInfos = append(typeInfos, typeInfo)
		}
		file.Close()
	}
	return typeInfos
}
