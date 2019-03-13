### 项目介绍
该项目是关于视频流数据的接收和分发器。主要是将接收到的视频流数据分发到其他节点或其他应用中。分发方式主要通过WEBSOCKET、TCP、UDP。如果在不要求实时性的情况下可以进行无限扩展。分发的视频流数据可以直接通过浏览器或者支持的播放器连接该服务进行播放。
服务采用配置的方式进行部署，以配置的数据进行数据流分发，配置看文档最后。这个项目要成为一个完成的项目主要包含数据采集部分、服务实现部分、数据应用部分。

项目主要因视频数据产生的，但是可以应用于多种场景，不仅仅是视频数据，比如可以应用于实时聊天、实时交互方面。如果加入数据存储，可以应用的方面将会比较多的。

### 部署方式
#### 1、ffmpeg推流到单节点
![ffmpeg推流到单节点](https://github.com/shockerjue/tsstream/blob/master/img/bushu1.png)
该方式是本地使用ffmpeg将捕获的视频流数据通过网络(HTTP、TCP、UDP)将视频数据推送到该项目实现的服务中，之后该服务将接收到视频数据分发到与之连接的客户端。客户端使用相应的应用进行渲染或者是使用数据。

#### 2、ffmpeg推流到单节点，节点内部扩展
![ffmpeg推流到单节点，节点内部扩展](https://github.com/shockerjue/tsstream/blob/master/img/bushu2.png)

该方式是本地使用ffmpeg将捕获的视频流数据通过网络(HTTP、TCP、UDP)将视频数据推送到该项目实现的服务中。在远程主机节点中部署多个服务，将收到的视频数据分发到内部部署的服务中。由内部的服务将视频数据分发到与服务连接的客户应用中进行渲染。这种方式主要是为了在单节点中扩展并发量。

#### 3、ffmpeg推流到本地服务，本地服务器对流到单节点
![架构图](https://github.com/shockerjue/tsstream/blob/master/img/bushu3.png)
主要是在本地部署一个服务，将ffmpeg捕获的视频流数据通过网络(HTTP、TCP、UDP)推送到本地服务中。由本地服务将收到的ffmpeg流数据分发到远程的节点中，远程节点在对流进行相应的处理。可以在远程节点中进行内部服务扩展以增加并发量。当然最后始终是要将流分发到与服务连接的客户应中进行渲染。

#### 4、ffmpeg推流到本地服务，本地服务将数据推流到多个节点中
![架构图](https://github.com/shockerjue/tsstream/blob/master/img/bushu4.png)
主要是在本地部署一个服务，将ffmpeg捕获的视频流数据通过网络(HTTP、TCP、UDP)推到本地服务中。本地服务器将ffmpeg流数据分发到远程的多个节点中，远程节点在对流进行相应的处理。可以在远程节点中进行内部服务扩展以增加并发量。当然最后始终是要将流分发到与服务连接的客户应中进行渲染。

#### 5、ffmpeg推流到本地服务，本地服务将数据推流到多个节点中，各节点在进行推流扩展
![架构图](https://github.com/shockerjue/tsstream/blob/master/img/bushu5.png)
主要是将ffmpeg捕获的数据通过服务器进行无限的扩展，这种方式如果对实时性要求不是很高的话，可以无限进行扩张，直到扩张到网络边界。当然最后始终是要将流分发到与服务连接的客户应中进行渲染。

### 将部署的对应端口打开
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

while true
do
    # ffmpeg -threads 2 -rtbufsize 16M -f video4linux2   -framerate 20 -video_size 1280x720  -i /dev/video0 -f mpegts -codec:v mpeg1video -s 640x480 -b:v 450k -bufsize 1300k -bufsize 2800k -maxrate 1500k -bf 0 -r 20 -muxdelay 0.001  http://
    ./ffmpeg -threads 2 -rtbufsize 16M -f video4linux2   -framerate 20 -video_size 640x480  -i /dev/video0 -f mpegts -codec:v mpeg1video -s 640x480 -b:v 256k -bufsize 1300k -bufsize 2800k -maxrate 1500k -bf 0 -r 20 -muxdelay 0.001  tcp://
    # ./ffmpeg -threads 2 -rtbufsize 16M -f video4linux2   -framerate 20 -video_size 640x480  -i /dev/video0 -f mpegts -codec:v mpeg1video -s 640x480 -b:v 450k -bufsize 1300k -bufsize 2800k -maxrate 1500k -bf 0 -r 20 -muxdelay 0.001  udp://
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
```


### 测试说明
本框架提供了具体的测试配置说明，测试目录在demo目录下
```
    client      这个目录是客户端文件，这个目录主要是通过web服务的方式使用
    push        这个目录是用于推流的目录，里面包含了linux、mac、windows等平台的推流脚本
    server      这个目录包含的是框架的部署方式以及部署程序
```
具体的使用说明可以查看demo目录下的说明文件。


### 节点监控
很多时候，我们需要知道对应的节点是否是活动的。
![节点监控](https://github.com/shockerjue/tsstream/blob/master/img/monitor.jpg)