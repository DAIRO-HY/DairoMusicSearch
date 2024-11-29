package set

import (
	"DairoMusicSearch/WebHandle"
	"net/http"
	"os"
)

// 路由设置
func Init() {
	http.HandleFunc("/set/cookie", WebHandle.ApiHandler(cookie))
}

func cookie(request *http.Request) string {
	//解析post表单
	err := request.ParseForm()
	if err != nil {
		return err.Error()
	}
	postForm := request.PostForm

	//得到cookie内容
	cookieTxt := postForm.Get("cookie")
	file, err := os.OpenFile("./youtube-cookie.txt", os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return err.Error()
	}
	if _, err := file.WriteString(cookieTxt); err != nil {
		return err.Error()
	}
	return "上传成功"
}
