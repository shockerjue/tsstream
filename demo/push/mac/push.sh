#!/bin/bash

# while true
# do
    ./ffmpeg -threads 8 -rtbufsize 16M  -framerate 20 -pixel_format yuyv422 -video_size  640x480  -f avfoundation -i "1" -f mpegts  -codec:v mpeg1video -s 640x480 -b:v 128k -bufsize 128k -bufsize 512k -maxrate 1024k -bf 0 -r 20 -muxdelay 0.001  tcp://120.78.156.129:50001
    ping -c 3 127.0.0.1 > /dev/null
# done

exit
