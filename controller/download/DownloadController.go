package download

import (
	"DairoMusicSearch/WebHandle"
	"DairoMusicSearch/extension/Number"
	"DairoMusicSearch/util/DownloadUtil"
	"DairoMusicSearch/util/YoutubeUtil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// 路由设置
func Init() {
	http.HandleFunc("/d/collect_info", WebHandle.ApiHandler(musicCollectInfo))
	http.HandleFunc("/d/music", WebHandle.ApiHandler(music))
	http.HandleFunc("/d/lrc", WebHandle.ApiHandler(lrc))
	http.HandleFunc("/d/proxy", WebHandle.ApiHandler(proxy))
}

func test(writer http.ResponseWriter, request *http.Request) {

	//range123 := request.Header.Get("range")
	//log.Println(range123)
	//
	//// 设置 Content-Type 头部信息
	//writer.Header().Set("Content-Type", "text/plain;charset=UTF-8")
	//writer.WriteHeader(http.StatusOK) // 设置状态码
	//writer.Write([]byte("SUCCESS"))

	//download("C:\\develop\\ideaIU-2024.1.2.win.zip", writer, request)
}

//    /**
//     * 下载音乐
//     * @param videoId 视频ID
//     * @param quality 音质 单位比特率
//     */
//    @CrossOrigin(origins = ["*"])
//    @GetMapping("/music")
//    fun music(request: HttpServletRequest, response: HttpServletResponse, videoId: String, quality: Int? = null) {
//        response.resetBuffer()
//
//        val os = System.getProperty("os.name").lowercase(Locale.getDefault())
//        if (os.contains("win")) {//测试用
//            response.status = HttpStatus.OK.value()
//            response.addHeader("Content-Type", "audio/mp3")
//            DownloadController::class.java.classLoader.getResource("test.mp3").openStream().use {
//                it.transferTo(response.outputStream)
//            }
//            return
//        }
//
//        val range = request.getHeader("range")
//        if(!range.isNullOrEmpty() && range != "bytes=0-"){//不支持选择范围
//            response.status = HttpStatus.REQUESTED_RANGE_NOT_SATISFIABLE.value()
//            return
//        }
//        response.resetBuffer()
//
//        //客户端输出流
//        val oStream = response.outputStream
//        val error = YoutubeUtil.getMusic(videoId, quality) {
//            response.status = HttpStatus.OK.value()
//            response.addHeader("Content-Type", "audio/mp3")
//            it.transferTo(oStream)
//        }
//        if (error != null) {
//            response.status = HttpStatus.SERVICE_UNAVAILABLE.value()
//            response.addHeader("Content-Type", "text/plain")
//            oStream.write(error.toByteArray())
//        }
//    }

//    /**
//     * 下载音乐
//     * @param videoId 视频ID
//     */
//    @CrossOrigin(origins = ["*"])
//    @ResponseBody
//    @GetMapping("/music_info/{videoId}")
//    fun musicInfo(@PathVariable("videoId") videoId: String): AudioInfo? {
//
//        //得到mp3文件
//        val mp3File = this.youtube.getMusic(videoId)
//        val info = this.cacheFile.getMp3Info(mp3File)
//        return info
//    }

/**
 * 音乐采集情况并发起采集请求
 * @param videoId 视频ID
 */
//@GetMapping("/music/{videoId}/collect_info")
func musicCollectInfo(request *http.Request) string {
	rootFolder := YoutubeUtil.RootDir()
	query := request.URL.Query()

	//得到视频ID
	videoId := query.Get("videoId")

	//获取开始采集时间
	collectStartTime := YoutubeUtil.GetCollectStartTime(videoId)

	//记录已经采集时间
	collectTimer := ""

	//已经耗时
	if collectStartTime != 0 { //如果采集已经开始,这里计算出已经耗时
		collectTimer = Number.ToTimeFormat((time.Now().UnixMilli() - collectStartTime) / 1000)
	} else {
		collectTimer = "0S"
	}

	var downloadFile string
	var mp3File string

	//获取采集相关的文件
	musicFileList, _ := os.ReadDir(rootFolder) //  rootFolder.listFiles()?.filter { it.name.startsWith(videoId) }
	for _, entry := range musicFileList {
		if strings.HasPrefix(entry.Name(), videoId) {
			if strings.HasSuffix(entry.Name(), ".part") || strings.HasSuffix(entry.Name(), ".webm") { //这是一个正在下载中的文件
				downloadFile = rootFolder + "/" + entry.Name()
			} else if strings.HasSuffix(entry.Name(), ".mp3") { //正在转换的mp3
				mp3File = rootFolder + "/" + entry.Name()
			}
		}
	}
	if len(downloadFile) == 0 && len(mp3File) == 0 {

		//去采集
		go YoutubeUtil.RequestCollectMusic(videoId)
		downloadingThreadCount := YoutubeUtil.CollectingCount()
		if downloadingThreadCount == 0 {
			return "准备采集:" + collectTimer
		}
		return "准备采集" + collectTimer + ",当前采集数:" + strconv.Itoa(downloadingThreadCount)
	}
	if len(mp3File) != 0 { //mp3文件已经存在，说明正在转码中
		if collectStartTime == 0 { //文件转码已经完成
			return "OK"
		}
		info, _ := os.Stat(mp3File)
		size := Number.ToDataSize(info.Size())
		return "正在转码:" + size + "(" + collectTimer + ")"
	}
	if len(downloadFile) != 0 { //正在下载中的文件
		info, _ := os.Stat(downloadFile)
		size := Number.ToDataSize(info.Size())
		return "正在采集:" + size + "(" + collectTimer + ")"
	}
	return ""

	//if (musicFileList.isNullOrEmpty()) {
	//
	//    //请求开始采集
	//    this.youtube.requestCollectMusic(videoId)
	//    val downloadingThreadCount = this.youtube.collectingCount
	//    if (downloadingThreadCount == 0) {
	//        return "准备采集$collectTimer"
	//    }
	//    return "准备采集$collectTimer,当前采集数:${downloadingThreadCount}"
	//}

	////正在下载中的文件
	//val downloadingFile = musicFileList.find { it.name.endsWith(".part") }
	//if (downloadingFile != null) {
	//    return "正在采集:${downloadingFile.length().dataSize}($collectTimer)"
	//}
	//if (collectStartTime != null) {//采集中
	//
	//    // 转码中
	//    val mp3File = musicFileList.find { it.name.endsWith(".mp3") } ?: return "准备转码($collectTimer)"
	//    return "正在转码:${mp3File.length().dataSize}($collectTimer)"
	//}else{//采集已经完成
	//
	//    //正在下载中的文件
	//    val mp3File = musicFileList.find { it.name == "$videoId.mp3" }
	//    if (mp3File != null) {//转码完成
	//        return "OK"
	//    }
	//}

	//未知结果
	//return "UNKNOWN:" + musicFileList.joinToString { it.name }
}

/**
 * 下载音乐
 * @param videoId 视频ID
 */
func music(writer http.ResponseWriter, request *http.Request) string {
	query := request.URL.Query()

	//得到视频ID
	videoId := query.Get("videoId")
	if strings.HasSuffix(videoId, "/collect_info") { //这里兼容旧版本，这是一个获取采集信息的url
		videoId = strings.ReplaceAll(videoId, "/collect_info", "")
		rootFolder := YoutubeUtil.RootDir()

		//获取开始采集时间
		collectStartTime := YoutubeUtil.GetCollectStartTime(videoId)

		//记录已经采集时间
		collectTimer := ""

		//已经耗时
		if collectStartTime != 0 { //如果采集已经开始,这里计算出已经耗时
			collectTimer = Number.ToTimeFormat((time.Now().UnixMilli() - collectStartTime) / 1000)
		} else {
			collectTimer = "0S"
		}

		var downloadFile string
		var mp3File string

		//获取采集相关的文件
		musicFileList, _ := os.ReadDir(rootFolder) //  rootFolder.listFiles()?.filter { it.name.startsWith(videoId) }
		for _, entry := range musicFileList {
			if strings.HasPrefix(entry.Name(), videoId) {
				if strings.HasSuffix(entry.Name(), ".part") || strings.HasSuffix(entry.Name(), ".webm") { //这是一个正在下载中的文件
					downloadFile = rootFolder + "/" + entry.Name()
				} else if strings.HasSuffix(entry.Name(), ".mp3") { //正在转换的mp3
					mp3File = rootFolder + "/" + entry.Name()
				}
			}
		}
		if len(downloadFile) == 0 && len(mp3File) == 0 {

			//去采集
			go YoutubeUtil.RequestCollectMusic(videoId)
			downloadingThreadCount := YoutubeUtil.CollectingCount()
			if downloadingThreadCount == 0 {
				return "准备采集:" + collectTimer
			}
			return "准备采集" + collectTimer + ",当前采集数:" + strconv.Itoa(downloadingThreadCount)
		}
		if len(mp3File) != 0 { //mp3文件已经存在，说明正在转码中
			if collectStartTime == 0 { //文件转码已经完成
				return "OK"
			}
			info, _ := os.Stat(mp3File)
			size := Number.ToDataSize(info.Size())
			return "正在转码:" + size + "(" + collectTimer + ")"
		}
		if len(downloadFile) != 0 { //正在下载中的文件
			info, _ := os.Stat(downloadFile)
			size := Number.ToDataSize(info.Size())
			return "正在采集:" + size + "(" + collectTimer + ")"
		}
		return ""
	}

	//得到mp3文件
	mp3File := YoutubeUtil.GetMusicFile(videoId)
	DownloadUtil.Download(mp3File, writer, request)

	//这里无需返回，只是兼容旧版本
	return ""
}

//
//    /**
//     * 下载音乐
//     * @param videoId 视频ID
//     */
//    @CrossOrigin(origins = ["*"])
//    @GetMapping("/video/{videoId}")
//    fun video(request: HttpServletRequest, response: HttpServletResponse, @PathVariable("videoId") videoId: String, size: Int?) {
//        val range = request.getHeader("range")
//        if (!range.isNullOrEmpty() && range != "bytes=0-") {//不支持选择范围
//            response.status = HttpStatus.REQUESTED_RANGE_NOT_SATISFIABLE.value()
//            return
//        }
//        response.resetBuffer()
//
//        //客户端输出流
//        val oStream = response.outputStream
//        val error = this.youtube.getVideo(videoId, size) {
//            response.status = HttpStatus.OK.value()
//            response.addHeader("Content-Type", "video/webm")
//            it.transferTo(oStream)
//        }
//        if (error != null) {
//            response.status = HttpStatus.SERVICE_UNAVAILABLE.value()
//            response.addHeader("Content-Type", "text/plain")
//            oStream.write(error.toByteArray())
//        }
//    }
//

// 下载歌词
func lrc(request *http.Request) string {
	query := request.URL.Query()

	//得到视频ID
	videoId := query.Get("videoId")
	return YoutubeUtil.GetLRC(videoId)
}

// 代理下载图片
func proxy(writer http.ResponseWriter, request *http.Request) any {
	query := request.URL.Query()

	//得到视频ID
	targetUrl := query.Get("url")
	resp, err := http.Get(targetUrl)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	writer.Header().Set("Content-Type", resp.Header["Content-Type"][0])
	buf := make([]byte, 8*1024)
	for {
		n, readErr := resp.Body.Read(buf)
		if readErr != nil {
			break
		}
		_, writeErr := writer.Write(buf[:n])
		if writeErr != nil {
			break
		}
	}
	return nil
}
