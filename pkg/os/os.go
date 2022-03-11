package os

import (
	"fmt"
	"gokit/pkg/i18n"
	"runtime"
	"strconv"
	"strings"
)
import "github.com/matishsiao/goInfo"

var (
	I18nKey_Description = i18n.I18nKey{Key: "os_description", Common: "系统相关工具"}
	description         = i18n.GetLocaleVal(I18nKey_Description)
)

type Os struct {
}

func (os *Os) Description() string {
	return description
}
func (os *Os) Run(args []string) {
	fmt.Println(osInfo())
}
func getOs() string {
	return runtime.GOOS
}
func osInfo() string {
	info, _ := goInfo.GetInfo()
	osInfo := strings.Builder{}
	osInfo.WriteString("\r\n")
	osInfo.WriteString("Platform:" + info.Platform + "\r\n")
	osInfo.WriteString("Hostname:" + info.Hostname + "\r\n")
	osInfo.WriteString("Os:" + info.OS + "\r\n")
	osInfo.WriteString("Kernel:" + info.Kernel + "\r\n")
	osInfo.WriteString("Core:" + info.Core + "\r\n")
	osInfo.WriteString("CPUS:" + strconv.Itoa(info.CPUs) + "\r\n")
	return osInfo.String()
}
