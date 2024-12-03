/**
 * 代码为自动生成，请勿手动修改
 */
package main

import (
	"DairoMusicSearch/controller/download"
	"DairoMusicSearch/controller/search"
	"DairoMusicSearch/controller/set"
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func init() {

	http.HandleFunc("/d/collect_info", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		videoId := query.Get("videoId")
		body := download.MusicCollectInfo(videoId)
		writeToResponse(writer, body)
	})
	http.HandleFunc("/d/music", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		videoId := query.Get("videoId")
		body := download.Music(writer, request, videoId)
		writeToResponse(writer, body)
	})
	http.HandleFunc("/d/lrc", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		videoId := query.Get("videoId")
		body := download.Lrc(videoId)
		writeToResponse(writer, body)
	})
	http.HandleFunc("/d/proxy", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		url := query.Get("url")
		body := download.Proxy(writer, url)
		writeToResponse(writer, body)
	})
	http.HandleFunc("/search/api", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		key := query.Get("key")
		body := search.Api(request, key)
		writeToResponse(writer, body)
	})
	http.HandleFunc("/set/cookie", func(writer http.ResponseWriter, request *http.Request) {
		body := set.Cookie(request)
		writeToResponse(writer, body)
	})
}

// 获取表单实例
func getForm[T any](request *http.Request) T {

	// 创建结构体实例
	targetForm := new(T)
	reflectForm := reflect.ValueOf(targetForm).Elem()
	argType := reflect.TypeOf(targetForm)
	//reflectForm := reflect.New(argType).Elem()
	query := request.URL.Query()

	//解析post表单
	request.ParseForm()
	postParams := request.PostForm

	//将参数转换成Map
	paramMap := make(map[string][]string)
	for key, v := range query {
		paramMap[strings.ToLower(key)] = v
	}
	for key, v := range postParams {
		paramMap[strings.ToLower(key)] = v
	}

	// 遍历结构体字段
	for j := 0; j < argType.NumField(); j++ {
		field := argType.Field(j)
		fieldName := field.Name

		//得到参数值
		value := paramMap[strings.ToLower(fieldName)]
		if value == nil {
			continue
		}

		// 设置字段值（这里我们设置为示例值）
		switch field.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

			// 设置整数字段
			intValue, _ := strconv.ParseInt(value[0], 10, 64)
			reflectForm.Field(j).SetInt(intValue)
		case reflect.String:
			reflectForm.Field(j).SetString(value[0]) // 设置字符串字段
		}
	}
	return *targetForm
}

// 返回结果
func writeToResponse(writer http.ResponseWriter, body any) {
	if body == nil {
		return
	}
	if body == "" {
		return
	}

	// 设置 Content-Type 头部信息
	writer.Header().Set("Content-Type", "text/plain;charset=UTF-8")

	switch returnBody := body.(type) {
	case string:
		writer.Write([]uint8(returnBody))
	case int:
		writer.Write([]uint8(strconv.Itoa(returnBody)))
	case int8:
		writer.Write([]uint8(strconv.Itoa(int(returnBody))))
	case int16:
		writer.Write([]uint8(strconv.Itoa(int(returnBody))))
	case int32:
		writer.Write([]uint8(strconv.Itoa(int(returnBody))))
	case int64:
		writer.Write([]uint8(strconv.FormatInt(returnBody, 10)))
	case error:
		// 设置 HTTP 状态码
		writer.WriteHeader(http.StatusInternalServerError) // 设置状态码
		jsonData, _ := json.Marshal(body)
		writer.Write(jsonData)
	default:
		jsonData, _ := json.Marshal(body)
		writer.Write(jsonData)
	}
}
