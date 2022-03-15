package server

import (
	"fmt"
	"gokit/pkg/common"
	sysHttp "net/http"
	"os"
	"strings"
)

func printUsage() {
	fmt.Println("使用方式 -bind[绑定地址端口] -local[本地路径]")
	fmt.Println("默认绑定地址localhost:8080,监听路径根目录,本地路径为当前目录")
}
func StartServer(args []string) {
	bind := "localhost:8080"
	local, _ := os.Getwd()
	for index, length := 0, len(args); length < len(args); index++ {
		arg := args[index]
		switch strings.ToLower(arg) {
		case "-bind":
			if common.LessIntThan(index+1, length) {
				bind = args[index+1]
				index++
			}
			break
		case "-local":
			if common.LessIntThan(index+1, length) {
				local = args[index+1]
				index++
			}
			break
		default:
		}
	}
	localHandler := sysHttp.FileServer(sysHttp.Dir(local))
	fmt.Println("服务已启动在", bind)
	sysHttp.ListenAndServe(bind, localHandler)
}
