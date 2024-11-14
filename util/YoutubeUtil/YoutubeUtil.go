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
 * 获取音乐
 * yt-dlp中，-o -  以流的形式输出时，--audio-format --audio-quality这样的参数无效输出的文件格式为opus
 * 要获取采集到的音频文件真正最好的比特率的mp3文件时使用用：
 *          yt-dlp --no-check-certificate -x --audio-format mp3 --audio-quality 0 "https://www.youtube.com/watch?v=d6f60-lxVQI"
 * @param videoId youtube视频ID
 * @param quality 文件比特率
 * @param soutCallback 数据流回调函数
 */
//func getMusic(videoId string, quality int, soutCallback: (iStream: InputStream) -> Unit): String? {
//    val quality = quality ?: 128
//
//    // FFMPEG对音频数据进行转换。具体操作说明：
//    // -i pipe:0：指定输入从管道读取。即yt-dlp输出流
//    // -f mp3：指定输出格式为 MP3。
//    // -b:a 128K：指定输出音频的比特率为 128Kbps。您可以根据需要调整比特率。
//    // -vn：禁用视频流输出。
//    // - 表示不输出文件,以流的形式输出
//
//    //这段指令在windows系统下报错，可能原始是管道操作符|导致的，有可能使用CMD.EXE能解决，尚未检证
//    val cmd =
//            "yt-dlp --no-check-certificate -x -o - \"https://www.youtube.com/watch?v=$videoId\" | ffmpeg -i pipe:0 -f mp3 -b:a ${quality}K -vn -"
//    return exec(cmd, soutCallback)
//}

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

	cmd := "yt-dlp --no-check-certificate -x --audio-format mp3 --audio-quality 0 -o \"" + mp3File + "\" https://www.youtube.com/watch?v=" + videoId + ""
	_, errResult, _ := ShellUtil.ExecToResult(cmd)
	lock.Lock()
	delete(collectVideoIdMap, videoId)
	lock.Unlock()
	if errResult != "" { //采集失败
		log.Println(errResult)
		return
	}
}

//
//    /**
//     * 下载Opus文件
//     * @param videoId youtube视频ID
//     * @param path 保存路径
//     */
//    private fun downloadOpus(videoId: String, path: String) {
//
//        //这段指令在windows系统下报错，可能原始是管道操作符|导致的，有可能使用CMD.EXE能解决，尚未检证
//        val cmd = """yt-dlp --no-check-certificate -x -o "${path}" "https://www.youtube.com/watch?v=$videoId""""
//        execShell(cmd)
//    }
//
//    /**
//     * 获取视频
//     * @param videoId youtube视频ID
//     * @param size 视频分辨率, 如720,代表视频高度为720
//     * @param soutCallback 数据流回调函数
//     */
//    fun getVideo(videoId: String, size: Int? = null, soutCallback: (iStream: InputStream) -> Unit): String? {
//
//        //视频分辨率参数
//        val res = if (size == null) {
//            ""
//        } else {
//            " -S \"+res:$size,codec,br\""
//        }
//        val cmd =
//                "yt-dlp --no-check-certificate$res -o - \"https://www.youtube.com/watch?v=$videoId\""
//        return exec(cmd, soutCallback)
//    }
//
/**
 * 获取字幕
 * @param videoId youtube视频ID
 */
func GetLRC(videoId string) string {
	folder := config.CacheFolder + "/lrc/" + videoId
	//try {
	command := "yt-dlp --no-check-certificate --write-subs --sub-lang all --skip-download https://www.youtube.com/watch?v=" + videoId + " -o \"" + folder + "/lrc\""
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

//
//    /**
//     * 执行指令
//     */
//    private fun exec(cmd: String, soutCallback: (iStream: InputStream) -> Unit): String? {
//        var process: Process? = null
//        try {
//            process = ShellUtil.getProcess(cmd)
//            val iStream = process.inputStream
//            thread {
//                //            iStream.bufferedReader().forEachLine {
//                //                println(it)
//                //            }
//
//                soutCallback(iStream)
//            }
//
//            //将错误流转换成字符串
//            val error = String(process.errorStream.readAllBytes())
//
//            //这里会线程阻塞
//            val code = process.waitFor()
//            if (code != 0) {//处理失败
//                return error
//            }
//            return null
//        } finally {
//
//            //强行终止子进程ß
//            process?.destroyForcibly()
//        }
//    }
//
//    /**
//     * 执行指令
//     */
//    private fun execShell(cmd: String): String {
//        var process: Process? = null
//        try {
//            process = ShellUtil.getProcess(cmd)
//
//            //记录当前时间
//            val now = System.currentTimeMillis()
//
//            //最大等待时间
//            val waitMaxTime = 2 * 60 * 1000
//            var isFinish = false
//            thread {
//                while (true) {
//                    if (isFinish) {
//                        return@thread
//                    }
//
//                    //已经执行时间
//                    val excutedTime = System.currentTimeMillis() - now
//                    if (excutedTime > waitMaxTime) {//如果执行时间超过5分钟,强行停止
//                        break
//                    }
//                    Thread.sleep(500)
////                    println("Shell已执行${excutedTime / 1000}秒,指令:$cmd")
//                }
//                process.destroyForcibly()
//            }
//
//            //这里会线程阻塞
//            val code = process.waitFor()
//            isFinish = true
//            if (code != 0) {//处理失败
////                if(!process.isAlive){
////                    throw RuntimeException("Shell执行失败,错误代码:$code  错误内容:执行超时")
////                }
//                val error = String(process.errorStream.readAllBytes())
//                throw RuntimeException("Shell执行失败,错误代码:$code  错误内容:" + error)
//            }
//            val rs = String(process.inputStream.readAllBytes())
//            return rs
//        } finally {
//
//            //强行终止子进程
//            process?.destroyForcibly()
//        }
//    }
//
//    companion object {
//        @JvmStatic
//        fun main(args: Array<String>) {
//            YoutubeUtil().exec("""yt-dlp --no-check-certificate -x --audio-format mp3 --audio-quality 0 -o "./12322313.mp3" "https://www.youtube.com/watch?v=ttnpG5vJpEY"""") {
//                val data = ByteArray(1024)
//                var len = 0
//                while (it.read(data).also { len = it } != -1) {
//                    println(String(data, 0, len))
//                }
//            }
//            println("sdfsd")
//        }
//    }
//}
