FROM ubuntu:22.04
MAINTAINER zhoulq<1064348912@qq.com>
RUN apt-get update
RUN yes|apt-get upgrade

WORKDIR /usr/local

#https://github.com/yt-dlp/yt-dlp
#####################################################################安装YT-DLP#####################################START
RUN yes|apt-get install python3
RUN ln -s /usr/bin/python3 /usr/bin/python
ADD https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp /usr/local/bin/yt-dlp
RUN chmod a+rx /usr/local/bin/yt-dlp
#####################################################################安装YT-DLP#######################################END


#YT-DLP依赖ffmpeg
RUN yes|apt-get -y install ffmpeg

RUN mkdir /app
WORKDIR /app

#下载程序
#ADD https://github.com/DAIRO-HY/DairoNPC/releases/download/1.0.0/dairo-npc-linux-amd64 /app
COPY dairo-music-search-linux-amd64 /app/dairo-music-search-linux-amd64

#授权可执行
RUN chmod +x dairo-music-search-linux-amd64

# 定义运行时的环境变量

#谷歌开发者KEY
ENV GOOGLE_API_KEY=xxxxxxxxxxxxx
#临时文件路径
ENV CACHE_FOLDER=./temp
#允许同时下载线程数
ENV MAX_DOWNLOAD_THREAD_COUNT=1
#WEB端口
ENV WEB_PORT=8021

#启动脚本
ENTRYPOINT /app/dairo-music-search-linux-amd64 -GoogleApiKey:$GOOGLE_API_KEY -CacheFolder:$CACHE_FOLDER -MaxDownloadThreadCount:$MAX_DOWNLOAD_THREAD_COUNT -WebPort:$WEB_PORT
