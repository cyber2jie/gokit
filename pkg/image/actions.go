package image

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/golang/freetype/truetype"
	"gokit/pkg/common"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type handler struct {
	encode func(io.Writer, image.Image) error
	decode func(io.Reader) (image.Image, error)
}
type handlerError struct {
}

func (h handlerError) Error() string {
	return "未找到handler"
}
func jpegEncode(w io.Writer, m image.Image) error {
	return jpeg.Encode(w, m, &jpeg.Options{Quality: 100})
}

var handlers = map[string]handler{
	"png":  handler{encode: png.Encode, decode: png.Decode},
	"jpeg": handler{encode: jpegEncode, decode: jpeg.Decode},
	"jpg":  handler{encode: jpegEncode, decode: jpeg.Decode},
}

func findHandler(path string) (handler, error) {
	lastDotIndex := strings.LastIndex(path, ".")
	if lastDotIndex != -1 {
		suffix := path[(lastDotIndex + 1):]
		suffix = strings.ToLower(suffix)
		handler := handlers[suffix]
		if &handler != nil {
			return handler, nil
		}
	}
	return handler{}, handlerError{}
}
func blur(args []string) {
	canBlur := true
	in, args := common.VisitOne(args)
	if in == "" {
		fmt.Println("邱少图片文件路径")
		canBlur = false
	}
	sigma, args := common.VisitOne(args)
	sigmaFloat, error := strconv.ParseFloat(sigma, 32)
	if error != nil {
		fmt.Println("sigma数错误")
		canBlur = false
	}
	out, _ := common.VisitOne(args)
	if out == "" {
		fmt.Println("缺少输出路径")
		canBlur = false
	}
	if canBlur {
		handler, error := findHandler(in)
		if error == nil {
			inFile, error := os.Open(in)
			if error == nil {
				outFile, error := os.OpenFile(out, os.O_CREATE|os.O_WRONLY, 0666)
				if error == nil {
					img, error := handler.decode(inFile)
					if error == nil {
						img = imaging.Blur(img, sigmaFloat)
						handler.encode(outFile, img)
						return
					}
				}
			}
		}
	}
	fmt.Println("使用方式：blur 图片文件路径 sigma数 输出路径")
}
func resise(args []string) {
	canResize := true
	in, args := common.VisitOne(args)
	if in == "" {
		fmt.Println("邱少图片文件路径")
		canResize = false
	}
	width, args := common.VisitOne(args)
	widthFloat, error := strconv.ParseFloat(width, 32)
	if error != nil {
		fmt.Println("width数错误")
		canResize = false
	}
	height, args := common.VisitOne(args)
	heightFloat, error := strconv.ParseFloat(height, 32)
	if error != nil {
		fmt.Println("height数错误")
		canResize = false
	}
	out, _ := common.VisitOne(args)
	if out == "" {
		fmt.Println("缺少输出路径")
		canResize = false
	}
	if canResize {
		handler, error := findHandler(in)
		if error == nil {
			inFile, error := os.Open(in)
			if error == nil {
				outFile, error := os.OpenFile(out, os.O_CREATE|os.O_WRONLY, 0666)
				if error == nil {
					img, error := handler.decode(inFile)
					if error == nil {
						img = imaging.Resize(img, int(widthFloat), int(heightFloat), imaging.Box)
						handler.encode(outFile, img)
						return
					}
				}
			}
		}
	}
	fmt.Println("使用方式：resize 图片文件路径 width height 输出路径")
}
func watermark(args []string) {
	canWaterMark := true
	in, args := common.VisitOne(args)
	if in == "" {
		fmt.Println("邱少图片文件路径")
		canWaterMark = false
	}
	fontPath, args := common.VisitOne(args)
	if fontPath == "" {
		fmt.Println("font路径错误")
		canWaterMark = false
	}
	text, args := common.VisitOne(args)
	if text == "" {
		fmt.Println("text错误")
		canWaterMark = false
	}
	fontSize, args := common.VisitOne(args)
	fontSizeFloat, error := strconv.ParseFloat(fontSize, 64)
	if error != nil {
		fmt.Println("fontSize错误")
		canWaterMark = false
	}
	x, args := common.VisitOne(args)
	xInt, error := strconv.ParseInt(x, 0, 64)
	if error != nil {
		fmt.Println("x位置错误")
		canWaterMark = false
	}
	y, args := common.VisitOne(args)
	yInt, error := strconv.ParseInt(y, 0, 64)
	if error != nil {
		fmt.Println("y位置错误")
		canWaterMark = false
	}
	out, _ := common.VisitOne(args)
	if out == "" {
		fmt.Println("缺少输出路径")
		canWaterMark = false
	}
	if canWaterMark {
		handler, error := findHandler(in)
		if error == nil {
			inFile, error := os.Open(in)
			if error == nil {
				outFile, error := os.OpenFile(out, os.O_CREATE|os.O_WRONLY, 0666)
				if error == nil {
					img, error := handler.decode(inFile)
					if error == nil {
						fontBytes, error := ioutil.ReadFile(fontPath)
						rgba := image.NewRGBA(image.Rect(0, 0, img.Bounds().Size().X, img.Bounds().Size().Y))
						draw.Draw(rgba, rgba.Bounds(), img, image.ZP, draw.Src)
						if error == nil {
							f, error := truetype.Parse(fontBytes)
							if error == nil {
								d := &font.Drawer{
									Dst: rgba,
									Src: image.Black,
									Face: truetype.NewFace(f, &truetype.Options{
										Size:    fontSizeFloat,
										DPI:     0,
										Hinting: font.HintingFull,
									}),
								}
								d.Dot = fixed.Point26_6{
									X: fixed.I(int(xInt)),
									Y: fixed.I(int(yInt)),
								}
								d.DrawString(text)
							}
						}
						handler.encode(outFile, rgba)
						return
					}
				}
			}
		}
	}
	fmt.Println("使用方式：watermark 图片文件路径 ttf文件路径 文本 fontSize x坐标 y坐标 输出路径")
}
