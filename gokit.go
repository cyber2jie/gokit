package gokit

import (
	"flag"
	"fmt"
	"gokit/pkg/image"
	"gokit/pkg/net"
	"gokit/pkg/os"
	"strings"
)

type Kit interface {
	Description() string
	Run(args []string)
}

const (
	kitOS    = "os"
	kitImage = "image"
	kitNet   = "net"
)

var kits = map[string]Kit{
	kitOS:    &os.Os{},
	kitImage: &image.Image{},
	kitNet:   &net.Net{},
}

func GoKit() {
	flag.Parse()
	kitKey := strings.ToLower(flag.Arg(0))
	kit := kits[kitKey]
	if kit == nil {
		printUsage(kitKey)
		return
	}
	kitArgs := flag.Args()[1:]
	kit.Run(kitArgs)
}
func printUsage(kitKey string) {
	if kitKey != "" {
		fmt.Println(fmt.Sprintf("未找到工具%s", kitKey))
	}
	fmt.Println("使用方式 gokit kit（工具） args (参数)")
	kitsBuilder := strings.Builder{}
	kitsBuilder.WriteString("可用kit（工具）")

	printKit(kitOS, &kitsBuilder)
	printKit(kitImage, &kitsBuilder)
	printKit(kitNet, &kitsBuilder)

	fmt.Println(kitsBuilder.String())
}
func printKit(kit string, kitsBuilder *strings.Builder) {
	toggle := "  "
	newLine := "\r\n"
	kitsBuilder.WriteString(newLine)
	kitsBuilder.WriteString(kit)
	kitsBuilder.WriteString(newLine)
	kitsBuilder.WriteString(toggle + kits[kit].Description())
	kitsBuilder.WriteString(newLine)
}
