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
ENV NPC_HOST=127.0.0.1
ENV NPC_KEY=xxxxxxxxx
ENV NPC_TCP=1781
ENV NPC_UDP=1782

#启动脚本
ENTRYPOINT /app/dairo-music-search-linux-amd64 -h:$NPC_HOST -k:$NPC_KEY -t:$NPC_TCP -u:$NPC_UDP
