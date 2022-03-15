package net

import (
	"fmt"
	"gokit/pkg/i18n"
	"gokit/pkg/net/http"
	"gokit/pkg/net/server"
	"os"
)

var (
	I18nKey_Description = i18n.I18nKey{Key: "net_description", Common: "网络相关工具"}
	description         = i18n.GetLocaleVal(I18nKey_Description)
	commander           = map[string]func(args []string){
		"httpClient":   httpClient,
		"simpleServer": simpleServer,
	}
)

func httpClient(args []string) {
	http.Run(args)
}

func simpleServer(args []string) {
	server.StartServer(args)
}

func commanders() []string {
	commanders := make([]string, len(commander))
	index := 0
	for cmd, _ := range commander {
		commanders[index] = cmd
		index++
	}
	return commanders
}
func printUsage(command string) {
	if command != "" {
		fmt.Println("未找到命令", command)
	}
	fmt.Println("使用方式，命令加参数")
	fmt.Println("可用命令：", commanders())
}

type Net struct {
}

func (net *Net) Description() string {
	return description
}
func (net *Net) Run(args []string) {
	command := ""
	if len(args) > 0 {
		command = args[0]
		args = args[1:]
	}
	cmdExist := false
	for _, cmd := range commanders() {
		if cmd == command {
			cmdExist = true
			break
		}
	}
	if !cmdExist {
		printUsage(command)
		os.Exit(-1)
	}
	commander[command](args)
}
