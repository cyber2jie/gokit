package common

import (
	"bytes"
	"strings"
)

func VisitOne(args []string) (string, []string) {
	if len(args) > 0 {
		return args[0], args[1:]
	}
	return "", []string{}
}
func LessIntThan(a int, b int) bool {
	return a < b
}
func StartWith(str string, startWith string) bool {
	for index, char := range str {
		if index >= len(startWith) {
			return true
		}
		for startIndex, startWithChar := range startWith {
			if index == startIndex {
				if char != startWithChar {
					return false
				}
			}
		}
	}
	return len(str) == len(startWith)
}
func Reverse(str string) string {
	runes := bytes.Runes([]byte(str))
	builder := strings.Builder{}
	for index := len(runes); index > 0; index-- {
		builder.WriteRune(runes[index-1])
	}
	return builder.String()
}
func EndWith(str string, endWith string) bool {
	return StartWith(Reverse(str), Reverse(endWith))
}
