FROM ubuntu:22.04
MAINTAINER zhoulq<1064348912@qq.com>
RUN apt-get update
RUN yes|apt-get upgrade

#更新CA证书，解决google api访问证书报错的BUG
RUN apt-get install -y ca-certificates
RUN update-ca-certificates

#再次更新软件元
RUN apt-get update

WORKDIR /usr/local

#####################################################################安装YT-DLP#####################################START
#https://github.com/yt-dlp/yt-dlp
RUN yes|apt-get install python3
RUN ln -s /usr/bin/python3 /usr/bin/python
ADD https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp /usr/local/bin/yt-dlp
RUN chmod a+rx /usr/local/bin/yt-dlp
#####################################################################安装YT-DLP#######################################END

#####################################################################安装FFMPEG#####################################START
#ADD https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz ffmpeg-release-amd64-static.tar.xz
ADD ffmpeg-7.0.2-amd64-static.tar.xz .
RUN ln -s /usr/local/ffmpeg-7.0.2-amd64-static/ffmpeg /usr/local/bin/ffmpeg
RUN ln -s /usr/local/ffmpeg-7.0.2-amd64-static/ffprobe /usr/local/bin/ffprobe
#####################################################################安装FFMPEG#####################################END


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
ENTRYPOINT /app/dairo-music-search-linux-amd64 GoogleApiKey:$GOOGLE_API_KEY CacheFolder:$CACHE_FOLDER MaxDownloadThreadCount:$MAX_DOWNLOAD_THREAD_COUNT WebPort:$WEB_PORT
