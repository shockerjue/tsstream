@echo off
setlocal enabledelayedexpansion
:main
ffmpeg -threads 2 -rtbufsize 16M -f dshow -framerate 20 -video_size 640x480 -i video="Logitech Webcam C930e" -f mpegts -codec:v mpeg1video -s 640x480 -b:v 1300k -bufsize 2800k -maxrate 1500k -bf 0 -r 20 -muxdelay 0.001  tcp://47.75.129.127:50001
ping -n 3 127.0.0.1>nul
goto main
pause