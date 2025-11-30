package YoutubeUtil

import (
	"DairoMusicSearch/config"
	"DairoMusicSearch/util/ShellUtil"
	"log"
	"os"
	"sync"
	"time"
)

// 正在采集中的ID对应的采集开始时间
var collectVideoIdMap = make(map[string]int64)

var lock sync.Mutex

/**
 * 正在下载中的线程数
 */
//    @Volatile
//    var downloadingThreadCount = 0

// 歌曲保存目录
func RootDir() string {
	return config.CacheFolder + "/music"
}

/**
 * 获取当前采集数
 */
func CollectingCount() int {
	return len(collectVideoIdMap)
}

/**
 * 判断是否正在采集中
 */
func GetCollectStartTime(videoId string) int64 {
	lock.Lock()
	startTime := collectVideoIdMap[videoId]
	lock.Unlock()
	return startTime
}

/**
 * 获取音乐文件
 */
func GetMusicFile(videoId string) string {
	return RootDir() + "/" + videoId + ".mp3"
}

/**
 * 请求采集音乐
 * yt-dlp中，-o -  以流的形式输出时，--audio-format --audio-quality这样的参数无效输出的文件格式为opus
 * 要获取采集到的音频文件真正最好的比特率的mp3文件时使用用：
 *          yt-dlp --no-check-certificate -x --audio-format mp3 --audio-quality 0 "https://www.youtube.com/watch?v=d6f60-lxVQI"
 * @param videoId youtube视频ID
 */
func RequestCollectMusic(videoId string) {
	lock.Lock()

	_, isExists := collectVideoIdMap[videoId]
	if isExists { //当前视频ID正在采集中
		lock.Unlock()
		log.Println(videoId + "正在采集中")
		return
	}

	if len(collectVideoIdMap) >= config.MaxDownloadThreadCount {
		//throw RuntimeException("服务器繁忙,请稍后再试")
		lock.Unlock()
		log.Println("达到最大采集并发数")
		return
	}

	//mp3文件
	mp3File := GetMusicFile(videoId)
	_, err := os.Stat(mp3File)
	if os.IsExist(err) { //指定mp3文件已经存在,说明已经采集完成
		lock.Unlock()
		return
	}
	collectVideoIdMap[videoId] = time.Now().UnixMilli()
	lock.Unlock()

	cmd := `yt-dlp --cookies ./data/youtube-cookie.txt --no-check-certificate -x --audio-format mp3 --audio-quality 0 -o "` + mp3File + `" "https://www.youtube.com/watch?v=` + videoId + `"`
	_, errResult, _ := ShellUtil.ExecToResult(cmd)
	lock.Lock()
	delete(collectVideoIdMap, videoId)
	lock.Unlock()
	if errResult != "" { //采集失败
		log.Println(errResult)
		return
	}
}

/**
 * 获取字幕
 * @param videoId youtube视频ID
 */
func GetLRC(videoId string) string {
	folder := config.CacheFolder + "/lrc/" + videoId
	command := `yt-dlp --cookies ./data/youtube-cookie.txt --no-check-certificate --write-subs --sub-lang all --skip-download "https://www.youtube.com/watch?v=` + videoId + `" -o "` + folder + `/lrc"`
	_, errResult, _ := ShellUtil.ExecToResult(command)
	if errResult != "" { //采集失败
		log.Println(errResult)
		return ""
	}
	var lrc = ""

	//多个歌词分割线
	const pageSplitStr = "\n<----------->\n"
	entries, _ := os.ReadDir(folder)
	for _, entry := range entries {
		content, _ := os.ReadFile(folder + "/" + entry.Name())
		lrc += string(content) + pageSplitStr
	}

	//删除文件夹
	os.RemoveAll(folder)
	if len(lrc) != 0 { //去掉最后的分割线
		lrc = lrc[0 : len(lrc)-len(pageSplitStr)]
	}
	return lrc
}
