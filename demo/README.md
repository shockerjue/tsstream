# 部署测试
### 部署框架图
![ffmpeg推流到单节点，节点内部扩展](https://github.com/shockerjue/tsstream/blob/master/img/demo.png)
上图是一个单节点部署架构。

## 服务器部署
### normal配置
```
# app 的配置
[app]
# 允许模式,normal(正常模式，单进程方式),extra(扩展方式，这里多数是和配置有关)
RunMode = normal

# 流接收服务配置(HTTP/UDP/TCP)[normal]
# 接收来自ffmpeg数据服务
[backstream]
Protocol = TCP
Bind = 0.0.0.0
Port = 50001
# 接受流的秘钥或者密码
Secret = key_576588

# 主要用于识别创建接收通过UDP/TCP推送过来的数据的服务类型以及接收流服务信息[extra]
[pushstream]
Protocol = UDP
Bind = 0.0.0.0
Port = 55002


# 分发流配置信息，注是将流分发到CUSTOM，或者TCP/UDP服务中，是normal和extra共用的配置
# Protocol设置为CUSTOM，则分发到当前进程的CUSTOM服务
# 如果Protocol设置为TCP/UDP，则可以分发到其他进程，或者其他机器中
[dispatch]
# 分发使用的协议 UDP/TCP/CUSTOM(注：CUSTOM协议是分发到配置中custom配置中服务中,是直接与用户连接的服务中)
Protocol = UDP

# 分发的UDP和TCP的连接地址以及端口列表,只有protocol设置为UDP/TCP才会被使用
Hosts = 127.0.0.1,127.0.0.1
Ports = 55001,55002


# 定义接受数据的ip地址列表,如果没有值，则可以接受任意IP的数据，如果设置了IPS,则只有在列表中中发来的数据才能被接受
[authority]
# 是否启用IP授权验证
Use = False
Ips = 192.168.0.107,192.168.12.7,192.168.0.114


# 定义用户节点的服务信息,如果dispatch的分发协议设置为CUSTOM，则会根据下面的地址信息创建对应的服务
[custom]
# 分发协议类型(WEBSOCKET/TCP)
Protocol = WEBSOCKET

# 如果使用TCP需要过滤的IP地址
FilterIps = 0.0.0.0

Bind = 0.0.0.0
Port = 50002

```

### extra1配置
```
# app 的配置
[app]
# 允许模式,normal(正常模式，单进程方式),extra(扩展方式，这里多数是和配置有关)
RunMode = extra

# 流接收服务配置(HTTP/UDP/TCP)[normal]
# 接收来自ffmpeg数据服务
[backstream]
Protocol = TCP
Bind = 0.0.0.0
Port = 50001
# 接受流的秘钥或者密码
Secret = key_576588

# 主要用于识别创建接收通过UDP/TCP推送过来的数据的服务类型以及接收流服务信息[extra]
[pushstream]
Protocol = UDP
Bind = 0.0.0.0
Port = 55001

# 分发流配置信息，注是将流分发到CUSTOM，或者TCP/UDP服务中，是normal和extra共用的配置
# Protocol设置为CUSTOM，则分发到当前进程的CUSTOM服务
# 如果Protocol设置为TCP/UDP，则可以分发到其他进程，或者其他机器中
[dispatch]
# 分发使用的协议 UDP/TCP/CUSTOM(注：CUSTOM协议是分发到配置中custom配置中服务中,是直接与用户连接的服务中)
Protocol = CUSTOM

# 分发的UDP和TCP的连接地址以及端口列表
Hosts = 127.0.0.1,127.0.0.1
Ports = 55001,55002


# 定义接受数据的ip地址列表,如果没有值，则可以接受任意IP的数据
[authority]
Use = False
Ips = 192.168.0.107,192.168.12.7,192.168.0.114


# 服务客户的配置信息，即直接服务器客户的配置
[custom]
# 分发协议类型(WEBSOCKET/TCP)
Protocol = WEBSOCKET

# 如果使用TCP需要过滤的IP地址，也就是推流过来的地址
FilterIps = 0.0.0.0
Bind = 0.0.0.0
Port = 56001

```

### extra2配置
```
# app 的配置
[app]
# 允许模式,normal(正常模式，单进程方式),extra(扩展方式，这里多数是和配置有关)
RunMode = extra

# 流接收服务配置(HTTP/UDP/TCP)[normal]
# 接收来自ffmpeg数据服务
[backstream]
Protocol = TCP
Bind = 0.0.0.0
Port = 50001
# 接受流的秘钥或者密码
Secret = key_

# 主要用于识别创建接收通过UDP/TCP推送过来的数据的服务类型以及接收流服务信息[extra]
[pushstream]
Protocol = UDP
Bind = 0.0.0.0
Port = 55002

# 分发流配置信息，注是将流分发到CUSTOM，或者TCP/UDP服务中，是normal和extra共用的配置
# Protocol设置为CUSTOM，则分发到当前进程的CUSTOM服务
# 如果Protocol设置为TCP/UDP，则可以分发到其他进程，或者其他机器中
[dispatch]
# 分发使用的协议 UDP/TCP/CUSTOM(注：CUSTOM协议是分发到配置中custom配置中服务中,是直接与用户连接的服务中)
Protocol = CUSTOM

# 分发的UDP和TCP的连接地址以及端口列表
Hosts = 127.0.0.1,127.0.0.1
Ports = 55001,55002


# 定义接受数据的ip地址列表,如果没有值，则可以接受任意IP的数据
[authority]
Use = False
Ips = 192.168.0.107,192.168.12.7,192.168.0.114


# 服务客户的配置信息，即直接服务器客户的配置
[custom]
# 分发协议类型(WEBSOCKET/TCP)
Protocol = WEBSOCKET

# 如果使用TCP需要过滤的IP地址，也就是推流过来的地址
FilterIps = 0.0.0.0
Bind = 0.0.0.0
Port = 56002

```

### 部署服务
将框架拷贝到需要部署的目录下，然后执行run.sh脚本，以允许对应的服务。
允许脚本以后，服务就启动完成，这个时候还需要配置具体的客户端访问服务，需要使用nginx代理来部署客户端访问具体的服务。
```
    upstream wstsstream {
        server 127.0.0.1:56001;
        server 127.0.0.1:56002;
    }

    server {
        listen 80;
        server_name  domain;


        location / {
                proxy_pass  http://wstsstream;
                proxy_http_version 1.1;
                proxy_set_header Upgrade $http_upgrade;
                proxy_set_header Connection $connection_upgrade;
        }
    }
```

## 客户端部署
客户端主要是通过网页的形式进行部署，部署之后通过浏览器进行访问。在部署之后修改播放地址为部署的服务器的地址。

```
<script type="text/javascript" src="jsmpeg.min.js"></script>
<script type="text/javascript">
    var canvas = document.getElementById('video-canvas');
    var url = 'wss://....'; // 播放地址//
    var player = new JSMpeg.Player(url,{canvas: canvas});
</script>
```

播放主要是引入jsmpeg.min.js文件，使用其中的JSMpeg.Player进行播放。

##### 部署
```
server {
    listen 80;
    root rootdir; # 客户端根目录

    index index.html index.htm index.nginx-debian.html index.php;

    server_name domain; # 部署的域名

    location / {
            try_files $uri $uri/ =404;
    }

    location ~ .*\.(js|css)?$ {
            add_header      Cache-Control no-cache;
            add_header      Cache-Control private;
            expires         0;
    }
}
```
部署主要采用的Nginx进行部署，其中通过浏览器直接进行访问。这样客户就可以直接使用浏览器进行访问具体的服务。


## 推送视频流数据
推送视频流时间使用ffmpeg自带的功能直接进推送。根据需要修改推送脚本，其中主要的是修改推送协议和地址。
```
-f mpegts -codec:v mpeg1video -s 640x360 -b:v 450k -bufsize 1300k -bufsize 2800k -maxrate 1500k -bf 0 -r 20 -muxdelay 0.001  tcp://x.x.x.x:50001
```
主要修改上面信息中的tcp://x.x.x.x，可以修改为http地址或者是udp地址,配置的信息就是normal配置的接收流的地址信息.