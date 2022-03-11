package http

import (
	"bytes"
	"fmt"
	"gokit/pkg/common"
	"io"
	"mime/multipart"
	sysHttp "net/http"
	"os"
	"strings"
)

type param struct {
	Name      string
	Value     string
	ValueType string
}

type paramError struct {
}

func (pe paramError) Error() string {
	return "参数格式错误"
}
func printUsage() {
	fmt.Println("使用方式 get/post[delete,put,option...] url 参数")
	fmt.Println(" -o filePath 输出位置")
	fmt.Println("参数形式 name=xxxx:类型(默认string,#:转义:号)")
	fmt.Println("类型类别：string,file,header")
}

func Run(args []string) {
	if len(args) == 0 {
		printUsage()
		return
	}
	method, args := common.VisitOne(args)
	if method == "" {
		printUsage()
		return
	}
	url, args := common.VisitOne(args)
	if url == "" {
		printUsage()
		return
	}
	writer := os.Stdout
	params := []param{}
	headers := map[string]string{}
	for index, length := 0, len(args); index < length; index++ {
		pm := args[index]
		switch strings.ToLower(pm) {
		case "-o":
			if len(args) > (index + 1) {
				path := args[index+1]
				outWriter, error := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
				if error != nil {
					fmt.Println("无法输出到路径:", error)
					return
				}
				writer = outWriter
				index++
			}
			continue
		default:
			//none
		}
		param, error := parseParam(pm)
		if error == nil {
			if param.ValueType == "header" {
				headers[param.Name] = param.Value
			} else {
				params = append(params, param)
			}
		}
	}

	switch strings.ToLower(method) {
	case "get":
		get(url, writer, params, headers)
		break
	case "post":
		post(url, writer, params, headers)
		break
	default:
		request, error := request(strings.ToUpper(method), url, params, headers)
		if error == nil {
			resp, error := sysHttp.DefaultClient.Do(request)
			if error == nil {
				afterRequest(resp, writer)
			}
		} else {
			fmt.Println("请求错误：", error)
		}
	}
}
func parseParam(pm string) (param, error) {
	pM := strings.Split(pm, "=")
	if len(pM) != 2 {
		return param{}, nil
	}
	mType := pM[1]
	value := strings.Builder{}
	pType := "string"
	for index := 0; index < len(mType); index++ {
		char := mType[index]
		if char == ':' {
			if len(mType) > (index + 1) {
				charNext := mType[(index + 1)]
				if charNext != ':' {
					_pType := mType[(index + 1):len(mType)]
					switch strings.ToLower(_pType) {
					case "file":
						pType = "file"
						break
					case "header":
						pType = "header"
						break
					default:
						pType = "string"
					}
					break
				} else {
					value.WriteByte(char)
					index++
					continue
				}
			}
			continue
		}
		value.WriteByte(char)
	}

	return param{Name: pM[0], Value: value.String(), ValueType: pType}, nil

}
func afterRequest(resp *sysHttp.Response, writer io.WriteCloser) {
	defer resp.Body.Close()
	defer writer.Close()
	buffLen := 1024 * 16
	buff := make([]byte, buffLen)
	for {
		readLen, error := resp.Body.Read(buff)
		if error == nil {
			if readLen == buffLen {
				writer.Write(buff)
			} else {
				writer.Write(buff[:readLen])
			}
		} else if error == io.EOF {
			writer.Write(buff[:readLen])
			break
		} else {
			break
		}
	}
}
func get(url string, writer io.WriteCloser, params []param, headers map[string]string) {
	request, error := request("GET", url, params, headers)
	if error == nil {
		resp, error := sysHttp.DefaultClient.Do(request)
		if error == nil {
			afterRequest(resp, writer)
		}
	}
}
func makeParamReader(params []param) (io.Reader, string) {
	writer := &bytes.Buffer{}
	multiWriter := multipart.NewWriter(writer)
	for _, param := range params {
		switch param.ValueType {
		case "string":
			multiWriter.WriteField(param.Name, param.Value)
			break
		case "file":
			fileWriter, error := multiWriter.CreateFormFile(param.Name, param.Value)
			if error == nil {
				file, error := os.Open(param.Value)
				defer file.Close()
				if error == nil {
					io.Copy(fileWriter, file)
				}
			}
			break
		default:
		}
	}
	multiWriter.Close()
	return writer, multiWriter.FormDataContentType()
}
func post(url string, writer io.WriteCloser, params []param, headers map[string]string) {
	request, error := request("POST", url, params, headers)
	if error == nil {
		resp, error := sysHttp.DefaultClient.Do(request)
		if error == nil {
			afterRequest(resp, writer)
		}
	}
}

func request(method string, url string, params []param, headers map[string]string) (*sysHttp.Request, error) {
	reader, contentType := makeParamReader(params)
	headers["Content-Type"] = contentType
	request, error := sysHttp.NewRequest(strings.ToUpper(method), url, reader)
	for key, value := range headers {
		request.Header.Add(key, value)
	}
	if error != nil {
		return nil, error
	}
	return request, nil
}
