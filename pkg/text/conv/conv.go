package conv

import (
	"fmt"
	"gokit/pkg/common"
	"gokit/pkg/common/converter"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"
)

func printUsage() {
	toggle := "  "
	fmt.Println("使用方式，转换器 参数")
	fmt.Println("可用转化器toLower,toUpper,replace,split,repeat")
	fmt.Println("通用参数:")
	fmt.Println(toggle + "-in 输入参数")
	fmt.Println(toggle + "-inType 输入参数类型[file]")
	fmt.Println(toggle + "--outPath 输出目录")
}
func Run(args []string) {
	converter, args := common.VisitOne(args)
	argParam := parse(args)
	switch strings.ToLower(converter) {
	case "tolower":
		toLower(argParam)
		break
	case "toupper":
		toUpper(argParam)
		break
	case "replace":
		replace(argParam)
		break
	case "split":
		split(argParam)
		break
	case "repeat":
		repeat(argParam)
		break
	default:
		printUsage()
	}
}

type argParam struct {
	in        string
	inType    string
	outPath   string
	otherArgs map[string]string
}

func (ap argParam) getArg(key string) string {
	return ap.otherArgs[key]
}
func parse(args []string) argParam {
	in, inType, outPath := "", "", ""
	otherArgs := map[string]string{}
	length := len(args)
	continu := false
	for index, arg := range args {
		if continu {
			continu = false
			continue
		}
		switch arg {
		case "-in":
			if common.LessIntThan(index+1, length) {
				in = args[index+1]
				continu = true
			}
			break
		case "-inType":
			if common.LessIntThan(index+1, length) {
				inType = args[index+1]
				continu = true
			}
			break
		case "-outPath":
			if common.LessIntThan(index+1, length) {
				outPath = args[index+1]
				continu = true
			}
			break
		default:
			if common.StartWith(arg, "-") {
				argName := arg[1:]
				if common.LessIntThan(index+1, length) {
					argVal := args[index+1]
					otherArgs[argName] = argVal
					continu = true
				}
			}
		}
	}
	return argParam{
		in:        in,
		inType:    inType,
		outPath:   outPath,
		otherArgs: otherArgs,
	}
}
func toLower(param argParam) {
	if param.in == "" {
		fmt.Println("toLower缺少in参数")
		return
	}
	writer := getWriter(param)
	defer writer.Close()
	content := loadContent(param)
	writer.WriteString(strings.ToLower(content))
}
func toUpper(param argParam) {
	if param.in == "" {
		fmt.Println("toUpper缺少in参数")
		return
	}
	writer := getWriter(param)
	defer writer.Close()
	content := loadContent(param)
	writer.WriteString(strings.ToUpper(content))
}
func replace(param argParam) {
	if param.in == "" {
		fmt.Println("replace缺少in参数")
		return
	}
	from := param.otherArgs["from"]
	to := param.otherArgs["to"]
	fromType := param.otherArgs["fromType"]
	if from == "" || to == "" {
		fmt.Println("缺少参数from,to")
		return
	}
	writer := getWriter(param)
	defer writer.Close()
	content := loadContent(param)
	switch fromType {
	case "regex":
		pattern, error := regexp.Compile(from)
		if error != nil {
			fmt.Println(error)
			return
		}
		writer.WriteString(pattern.ReplaceAllString(content, to))
	default:
		writer.WriteString(strings.ReplaceAll(content, from, to))
	}
}
func split(param argParam) {
	if param.in == "" {
		fmt.Println("split缺少in参数")
		return
	}
	splitter := param.otherArgs["splitter"]
	splitterType := param.otherArgs["splitterType"]
	join := param.otherArgs["join"]
	if splitter == "" {
		fmt.Println("split缺少splitter参数")
		return
	}
	if join == "" {
		join = " "
	}
	writer := getWriter(param)
	defer writer.Close()
	content := loadContent(param)
	switch splitterType {
	case "regex":
		pattern, error := regexp.Compile(splitter)
		if error != nil {
			fmt.Println(error)
			return
		}
		splitArray := pattern.Split(content, -1)
		writer.WriteString(strings.Join(splitArray, join))
	default:
		splitArray := strings.Split(content, splitter)
		writer.WriteString(strings.Join(splitArray, join))
	}
}
func repeat(param argParam) {
	if param.in == "" {
		fmt.Println("repeat缺少in参数")
		return
	}
	countStr := param.otherArgs["count"]
	if countStr == "" {
		fmt.Println("repeat缺少count参数")
		return
	}
	count, error := converter.Convert(countStr, reflect.Int)
	if error != nil {
		fmt.Println("repeat参数count非整数数字")
		return
	}
	writer := getWriter(param)
	defer writer.Close()
	content := loadContent(param)
	writer.WriteString(strings.Repeat(content, int(count.(int64))))

}
func loadContent(param argParam) string {
	switch param.inType {
	case "file":
		content, error := ioutil.ReadFile(param.in)
		if error == nil {
			return string(content)
		}
	default:
	}
	return param.in
}
func getWriter(param argParam) *os.File {
	writer := os.Stdout
	if param.outPath != "" {
		outFile, error := os.OpenFile(param.outPath, os.O_WRONLY|os.O_CREATE, 0666)
		if error == nil {
			writer = outFile
		}
	}
	return writer
}
