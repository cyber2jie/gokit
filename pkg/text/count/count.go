package count

import (
	"fmt"
	"gokit/pkg/common"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	_kb = 1024
	_mb = _kb * _kb
	_gb = _kb * _mb
)

func printUsage() {
	toggle := "  "
	fmt.Println("使用方式: 可用参数 文件 统计字符")
	fmt.Println(toggle + "-l 统计行数")
	fmt.Println(toggle + "-c 统计字符")
	fmt.Println(toggle + "-b 统计字节量")
}
func Run(args []string) {
	countType, args := common.VisitOne(args)
	filePath, args := common.VisitOne(args)
	if filePath == "" {
		fmt.Println("缺少文件参数")
		printUsage()
		return
	}
	file, error := os.Open(filePath)
	if error != nil {
		fmt.Println(error)
		return
	}
	info, error := file.Stat()
	if error != nil {
		fmt.Println(error)
		return
	}
	if info.IsDir() {
		fmt.Println("can't count dir")
		return
	}
	defer file.Close()
	switch strings.ToLower(countType) {
	case "-l":
		lineCount := 1
		// \r\n  one line,\n\r two lines
		prev := byte(0)
		for {
			buf := make([]byte, 1024)
			n, error := file.Read(buf)
			for index := 0; index < n; index++ {
				b := buf[index]
				if b == '\n' {
					if prev != '\r' {
						lineCount++
					}
				}
				if b == '\r' {
					lineCount++
				}
				prev = b
			}
			if error != nil {
				if error == io.EOF {
					break
				}
				fmt.Println(error)
				return
			}
		}
		fmt.Println("count", lineCount, "lines")
	case "-c":
		countStr, _ := common.VisitOne(args)
		prev := make([]byte, 0)
		bufLen := 1024
		count := 0
		if countStr != "" {
			buf, error := ioutil.ReadAll(file)
			if error != nil {
				fmt.Println(error)
				return
			}
			count += strings.Count(string(buf), countStr)
		} else {
			for {
				buf := make([]byte, bufLen)
				n, error := file.Read(buf)
				prevLen := len(prev)
				rangeBuf := make([]byte, prevLen+n)
				for index, _ := range rangeBuf {
					if index < prevLen {
						rangeBuf[index] = prev[index]
					} else {
						rangeBuf[index] = buf[index-prevLen]
					}
				}
				runeCount := utf8.RuneCount(rangeBuf)
				for index := len(rangeBuf) - 1; index >= 0; index-- {
					if utf8.RuneCount(rangeBuf[:index]) == runeCount-1 {
						prev = rangeBuf[index:len(rangeBuf)]
						count += runeCount - 1
						break
					}
				}
				if error != nil {
					if error == io.EOF {
						break
					}
					fmt.Println(error)
					return
				}
			}
			count += utf8.RuneCount(prev)
		}
		fmt.Println("count", count)
	case "-b":
		byteCount := info.Size()
		u, _ := common.VisitOne(args)
		uStr := "byte"
		count := float64(byteCount)
		switch u {
		case "-kb":
			uStr = "kb"
			count /= _kb
		case "-mb":
			uStr = "mb"
			count /= _mb
		case "-gb":
			uStr = "gb"
			count /= _gb
		}
		if uStr == "byte" {
			fmt.Println("count", strconv.FormatFloat(count, 'f', 0, 64), uStr)
		} else {
			fmt.Println("count", strconv.FormatFloat(count, 'f', 2, 64), uStr)
		}
	default:
		printUsage()
	}
}
