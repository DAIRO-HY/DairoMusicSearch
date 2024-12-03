package WebHandle

import (
	"DairoMusicSearch/config"
	"DairoMusicSearch/util/LogUtil"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func renderTemplate(w http.ResponseWriter, tmpl string) {

	// 解析嵌入的模板
	t, err := template.ParseFS(templatesFiles,
		"templates/"+tmpl,
		//"templates/include/head.html",
		//"templates/include/top-bar.html",
		//"templates/include/data_size_chart.html",
		//"templates/include/speed_chart.html",
	)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	// 设置 Content-Type 头部信息
	w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	t.Execute(w, nil)
}

// 路由处理
func htmlHandler(writer http.ResponseWriter, request *http.Request) {
	paths := strings.Split(request.URL.Path, "/")
	htmlFile := paths[len(paths)-1]
	if len(htmlFile) == 0 {
		htmlFile = "index"
	}
	renderTemplate(writer, htmlFile+".html")
}

// 通过这种方式将静态资源文件打包进二进制文件
//
//go:embed static/*
var staticFiles embed.FS

//go:embed templates/*
var templatesFiles embed.FS

func Start() {

	// 处理静态文件
	//fs := http.FileServer(http.FS(staticFiles))
	//http.Handle("/static/", fs)
	//
	// 设置路由
	http.HandleFunc("/", htmlHandler)

	port := config.WebPort

	// 启动服务器
	LogUtil.Info(fmt.Sprintf("WEB管理后台端口 :%s", port))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		LogUtil.Error(fmt.Sprintf("WEB管理后台启动失败 :%q", err))
		log.Fatal(err)
	}
	fmt.Printf("WEB管理后台端口 :%s\n", port)
}
