package download

import (
	"DairoMusicSearch/extension/Number"
	"DairoMusicSearch/util/DownloadUtil"
	"DairoMusicSearch/util/YoutubeUtil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

/**
 * 音乐采集情况并发起采集请求
 * @param videoId 视频ID
 */
//GET:/d/collect_info
func MusicCollectInfo(videoId string) string {
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

/**
 * 下载音乐
 * @param videoId 视频ID
 */
//GET:/d/music
func Music(writer http.ResponseWriter, request *http.Request, videoId string) string {
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

// Lrc 下载歌词
// GET:/d/lrc
func Lrc(videoId string) string {
	return YoutubeUtil.GetLRC(videoId)
}

// Proxy 代理下载图片
// GET:/d/proxy
func Proxy(writer http.ResponseWriter, url string) any {
	resp, err := http.Get(url)
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
