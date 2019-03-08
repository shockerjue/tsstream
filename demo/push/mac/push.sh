#!/bin/bash

while true
do
    ./ffmpeg -threads 8 -rtbufsize 16M  -framerate 30 -video_size 1280x720  -f avfoundation -i "1" -f mpegts  -codec:v mpeg1video -s 1280x720 -b:v 256k -bufsize 1300k -bufsize 2800k -maxrate 1500k -bf 0 -r 30 -muxdelay 0.001  tcp://47.75.129.127:50001
    ping -c 3 127.0.0.1 > /dev/null
done

exit
