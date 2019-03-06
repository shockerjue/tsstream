### 项目介绍
该项目是一个数据接收和分发器。主要是将接收到的数据分发到其他节点或其他服务中。分发方式主要通过websocket、TCP、UDP。并且在不要实时性的情况下可以无限扩展。主要应用于视频数据的接收和分发器。如果是视频数据的话，可以直接部署了是要浏览器或者支持的播放器进行播放。

### 部署方式
#### ffmpeg推流到单节点
![ffmpeg推流到单节点](https://github.com/shockerjue/tsstream/blob/master/img/bushu1.png)
    该方式在本地使用ffmpeg将捕获的视频数据通过通道将视频数据推送到该项目实现的服务中，然后该服务接收到视频数据以后将视频时间分发到客户端。这里就是将视频数据分发浏览器中进行渲染。

#### ffmpeg推流到单节点，节点内部扩展
![ffmpeg推流到单节点，节点内部扩展](https://github.com/shockerjue/tsstream/blob/master/img/bushu2.png)
    该方式是本地使用ffmpeg将捕获的视频数据通过通道将视频数据推送到该项目实现的服务中。在节点中部署多个服务，将收到的视频数据分发到内部部署的服务中。在有内部的服务将视频数据分发到浏览器中进行渲染。这种方式主要是为了扩展并发量。

#### ffmpeg推流到本地服务，本地服务器对流到单节点
![架构图](https://github.com/shockerjue/tsstream/blob/master/img/bushu3.png)
    主要是在本地部署一个服务，之后将ffmpeg捕获的视频数据推到本地服务中。本地服务器收到ffmpeg推来的流数据以后，将流数据分发到远程的节点中，之后远程节点在流进行相应的处理。

#### ffmpeg推流到本地服务，本地服务将数据推流到多个节点中
![架构图](https://github.com/shockerjue/tsstream/blob/master/img/bushu4.png)
    主要是在本地部署一个服务，之后将ffmpeg捕获的视频数据推到本地服务中。本地服务器收到ffmpeg推来的流数据以后，将流数据分发到远程的多个节点中，之后每个远程节点在流进行相应的处理。

#### ffmpeg推流到本地服务，粉底服务将数据推流到多个节点中，各节点在进行推流扩展
![架构图](https://github.com/shockerjue/tsstream/blob/master/img/bushu5.png)
    主要是将ffmpeg捕获的数据通过服务器进行无限的扩展，这种方式如果对实时性要求不是很高的话，可以无限进行扩张，知道扩张到网络边界。

### 将端口打开
```
sudo iptables -I INPUT -p tcp --dport 50001 -j ACCEPT &&
sudo iptables -I INPUT -p udp --dport 50001 -j ACCEPT &&
sudo iptables -I INPUT -p tcp --dport 50002 -j ACCEPT &&
sudo iptables -I INPUT -p tcp --dport 8088 -j ACCEPT &&
sudo iptables -I INPUT -p tcp --dport 55002 -j ACCEPT 
```


### 推送视频脚本
```
#!/bin/bash

# while true
# do
    # ffmpeg -threads 2 -rtbufsize 16M -f video4linux2   -framerate 20 -video_size 1280x720  -i /dev/video0 -f mpegts -codec:v mpeg1video -s 640x480 -b:v 450k -bufsize 1300k -bufsize 2800k -maxrate 1500k -bf 0 -r 20 -muxdelay 0.001  http://
    ./ffmpeg -threads 2 -rtbufsize 16M -f video4linux2   -framerate 20 -video_size 640x480  -i /dev/video0 -f mpegts -codec:v mpeg1video -s 640x480 -b:v 256k -bufsize 1300k -bufsize 2800k -maxrate 1500k -bf 0 -r 20 -muxdelay 0.001  tcp://
    # ./ffmpeg -threads 2 -rtbufsize 16M -f video4linux2   -framerate 20 -video_size 640x480  -i /dev/video0 -f mpegts -codec:v mpeg1video -s 640x480 -b:v 450k -bufsize 1300k -bufsize 2800k -maxrate 1500k -bf 0 -r 20 -muxdelay 0.001  udp://
    ping -c 3 127.0.0.1 > /dev/null
# done

exit
```

#### 推送视频参数设置
```
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
```

