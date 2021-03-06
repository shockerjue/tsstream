#!/bin/bash

while true
do
    ffmpeg -threads 2 -rtbufsize 16M -f video4linux2   -framerate 20 -video_size 1280x720  -i /dev/video0 -f mpegts -codec:v mpeg1video -s 640x360 -b:v 450k -bufsize 1300k -bufsize 2800k -maxrate 1500k -bf 0 -r 20 -muxdelay 0.001  tcp://x.x.x.x:50001
    ping -c 3 127.0.0.1 > /dev/null
done

exit

# -rtbufsize        设置最大的real-time帧内存使用，也就是视频流缓冲。等待发送的视频信息
# -threads          设置编码的线程数
# -framerate        设置输入输出的帧数 
# -video_size       设置输入视频的尺寸，这个值必须明确指定
# -f                设置输入输出格式
# -i                设置输入输出文件、设备、url等等
# -s                设置输出的视频帧尺寸
# -v:b              设置视频输出的最小比特率
# -bufsize          设置视频输出的缓冲大小  
# -maxrate          设置最大的输出比特率
# -r                设定帧速率
# -muxdelay         设置最大的解码延迟