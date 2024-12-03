package search

import (
	"DairoMusicSearch/config"
	"DairoMusicSearch/controller/search/form"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

/**
 * 音乐搜索API
 * @param key 搜索关键字
 */
//GET:/search/api
func Api(request *http.Request, key string) any {
	if len(key) == 0 {
		return []form.SearchForm{}
	}

	//将搜索关键字编码
	q := url.QueryEscape(key)

	//显示最大检索结果
	const limit = 30
	searchUrl :=
		"https://www.googleapis.com/youtube/v3/search?key=" + config.GoogleApiKey + "&type=video&part=snippet&q=" + q + "&maxResults=" + strconv.Itoa(limit)
	resp, err := http.Get(searchUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	formData := []form.SearchForm{}
	items := gjson.Get(string(data), "items").Array()

	scheme := "http"
	if request.TLS != nil {
		scheme = "https"
	}
	origin := scheme + "://" + request.Host
	for _, item := range items {
		logo := item.Get("snippet.thumbnails.high.url").Str

		//通过代理的方式下载图片
		logo = origin + "/d/proxy?url=" + url.QueryEscape(logo)
		formData = append(formData, form.SearchForm{
			Name:    item.Get("snippet.title").Str,
			Logo:    logo,
			VideoId: item.Get("id.videoId").Str,
		})
	}
	return formData
}
