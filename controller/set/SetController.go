package set

import (
	"net/http"
	"os"
)

// POST:/set/cookie
func Cookie(request *http.Request) string {

	//解析post表单
	err := request.ParseForm()
	if err != nil {
		return err.Error()
	}
	postForm := request.PostForm

	//得到cookie内容
	cookieTxt := postForm.Get("cookie")
	file, err := os.OpenFile("./data/youtube-cookie.txt", os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return err.Error()
	}
	if _, err := file.WriteString(cookieTxt); err != nil {
		return err.Error()
	}
	return "上传成功"
}
