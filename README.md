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


### 配置说明
```
# 配置模板文件
# app 的配置
[app]
# 运行模式,normal(正常模式:用于接收ffmpeg流数据),extra(扩展方式:用于接收由UDP/TCP分发过来的流数据并将其分发出去)
RunMode = 


# 接收来自ffmpeg数据服务配置
[backstream]
# 启动接收流的服务类型(HTTP/UDP/TCP),默认值为HTTP
Protocol = 

# 接收流服务器的地址配置信息
Bind = 0.0.0.0
Port = 

# 使用HTTP推流的时的秘钥或者密码
Secret = 


# 接收通过UDP/TCP推送数据的服务配置
[pushstream]
# 创建接收推送数据的服务类型（协议类型）[UDP/TCP]
Protocol = 

# 服务地址配置信息
Bind = 0.0.0.0
Port = 

# 分发流配置信息，注是将流分发到CUSTOM，或者TCP/UDP服务中，是normal和extra共用的配置
# Protocol设置为CUSTOM，则分发到当前进程的CUSTOM服务
# 如果Protocol设置为TCP/UDP，则可以分发到其他进程，或者其他机器中
[dispatch]
# 分发使用的协议 UDP/TCP/CUSTOM(注：CUSTOM协议是分发到配置中custom配置中服务中,是直接与用户连接的服务中)
Protocol = 

# 分发的UDP和TCP的连接地址以及端口列表,只有protocol设置为UDP/TCP才会被使用
# 如果Protocol设置为CUSTOM切custom中的Protocol设置为TCP，则下面应该保护custom中的地址信息
Hosts = 
Ports = 


# 定义授权的ip地址列表,如果值为空，则可以接受任意IP的数据
[authority]
# 是否使用授权认证
Use = 
# 需要验证的IP地址列表
Ips = 


# 服务客户的配置信息，即直接服务器客户的配置
[custom]
# 分发协议类型(WEBSOCKET/TCP)
Protocol = WEBSOCKET

# 如果使用TCP需要过滤的IP地址，也就是推流过来的地址
FilterIps = 127.0.0.1
Bind = 0.0.0.0
Port = 
```