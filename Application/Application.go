package Application

import (
	"DairoMusicSearch/config"
	"DairoMusicSearch/util/LogUtil"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// 程序的初始化操作
func Init() {
	fmt.Println("------------------------------------------------------------------------")
	for _, it := range os.Args {
		fmt.Println(it)
	}
	fmt.Println("------------------------------------------------------------------------")
	LogUtil.Info("项目启动成功")
	for _, it := range os.Args {
		paramArr := strings.Split(it, ":")
		switch paramArr[0] {
		case "GoogleApiKey":
			config.GoogleApiKey = paramArr[1]
		case "CacheFolder":
			config.CacheFolder = paramArr[1]
		case "MaxDownloadThreadCount":
			maxDownloadThreadCount, err := strconv.ParseInt(paramArr[1], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			config.MaxDownloadThreadCount = int(maxDownloadThreadCount)
		case "WebPort":
			config.WebPort = paramArr[1]
		}
	}
}
