package DownloadUtil

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// 路由设置
func Init() {
	http.HandleFunc("/", test)
}

func test(writer http.ResponseWriter, request *http.Request) {

	//range123 := request.Header.Get("range")
	//log.Println(range123)
	//
	//// 设置 Content-Type 头部信息
	//writer.Header().Set("Content-Type", "text/plain;charset=UTF-8")
	//writer.WriteHeader(http.StatusOK) // 设置状态码
	//writer.Write([]byte("SUCCESS"))

	Download("C:\\develop\\ideaIU-2024.1.2.win.zip", writer, request)
}

/**
* 文件下载
* @param mp3File 文件
* @param request 客户端请求
* @param response 往客户端返回内容
 */
func Download(mp3File string, writer http.ResponseWriter, request *http.Request) {
	fileInfo, err := os.Stat(mp3File)
	if os.IsNotExist(err) { //文件不存在
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	//文件大小
	size := fileInfo.Size()

	//指定读取部分数据头部标识
	ranges := request.Header.Get("range")
	var start int64
	var end int64
	if len(ranges) == 0 {
		start = 0
		end = size - 1
	} else {
		//range格式：bytes=10-30 或者 bytes=10-30
		rangeArr := strings.Split(strings.ToLower(ranges)[6:], "-")
		start, _ = strconv.ParseInt(rangeArr[0], 10, 64)
		if len(rangeArr[1]) == 0 { //到文件末尾
			end = size - 1
		} else {
			end, _ = strconv.ParseInt(rangeArr[1], 10, 64)
			if end > size-1 { //超出了文件大小范围
				end = size - 1
			}
		}
		if start > end {
			writer.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			return
		}
		if start >= size {
			writer.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			return
		}
		writer.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, size))
	}
	writer.Header().Set("Content-Type", "audio/mp3")
	writer.Header().Set("Content-Length", strconv.FormatInt(end-start+1, 10))

	//告诉客户端,服务器支持请求部分数据
	writer.Header().Set("Accept-Ranges", "bytes")

	file, err := os.Open(mp3File)
	defer file.Close()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	//http.ResponseWriter发送状态码之后，再设置头部信息将会不生效，所以发送状态码一定要等所有头部信息设置完成之后再发送
	if len(ranges) == 0 {
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusPartialContent)
	}

	//跳过前面部分数据
	file.Seek(start, io.SeekStart)
	data := make([]byte, 16*1024) // 缓冲字节数组
	var total = start
	for {

		//计算还需要的数据长度
		needReadLen := int(end - total + 1)
		n, readErr := file.Read(data)
		if readErr != nil {
			//if readErr != io.EOF { //如果不是文件读取完成标志,理论上，这里不会发生该异常
			//	writer.WriteHeader(http.StatusInternalServerError)
			//}
			break
		}
		total += int64(n)
		if needReadLen <= n { //还需要的数据长度小于本次读取到的数据长度
			writer.Write(data[:needReadLen])
			break
		} else {
			_, writeErr := writer.Write(data[:n])
			if writeErr != nil { //可能客户端已经关闭停止
				break
			}
		}
	}
}
