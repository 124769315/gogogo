# 作业3.2
## 创建Dockerfile
```
[root@master1 3.2]# cat Dockerfile 
FROM golang:latest
WORKDIR $GOPATH/src/httpserver
COPY httpserver.go $GOPATH/src/httpserver/
RUN go build $GOPATH/src/httpserver/httpserver.go
ENTRYPOINT ["./httpserver"]
[root@master1 3.2]# 
```

## 创建docker镜像
```
[root@master1 3.2]# docker build -t go-httpserver .
Sending build context to Docker daemon  3.584kB
Step 1/5 : FROM golang:latest
 ---> 8b86bf336a01
Step 2/5 : WORKDIR $GOPATH/src/httpserver
 ---> Using cache
 ---> 2b49c9cc33c2
Step 3/5 : COPY httpserver.go $GOPATH/src/httpserver/
 ---> Using cache
 ---> 1de7fa1ecc07
Step 4/5 : RUN go build $GOPATH/src/httpserver/httpserver.go
 ---> Running in fb633f4791c1
Removing intermediate container fb633f4791c1
 ---> fe46b9e61ceb
Step 5/5 : ENTRYPOINT ["./httpserver"]
 ---> Running in 0569840f0594
Removing intermediate container 0569840f0594
 ---> 13d8039b754f
Successfully built 13d8039b754f
Successfully tagged go-httpserver:latest
[root@master1 3.2]# 

```

## 检查docker镜像
##### docker image ls
```
[root@master1 3.2]# docker image ls
REPOSITORY      TAG       IMAGE ID       CREATED              SIZE
go-httpserver   latest    13d8039b754f   About a minute ago   947MB
```


## 启动http容器
```
[root@master1 3.2]# docker run -d go-httpserver 
4b2634582a819afcd47f739a487920d70baee4b31188f848dc3caaa75bb42a5c
[root@master1 3.2]# 
[root@master1 3.2]# 
[root@master1 3.2]# docker ps
CONTAINER ID   IMAGE           COMMAND          CREATED         STATUS         PORTS     NAMES
4b2634582a81   go-httpserver   "./httpserver"   4 seconds ago   Up 3 seconds             vigorous_bose
[root@master1 3.2]# 
```
## 查看容器ip配置
```
[root@master1 3.2]# docker inspect 4b2634582a81 | grep Pid
            "Pid": 27001,
            "PidMode": "",
            "PidsLimit": null,
[root@master1 3.2]# nsenter -t 27001 -n ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
24: eth0@if25: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
    link/ether 02:42:ac:11:00:02 brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 172.17.0.2/16 brd 172.17.255.255 scope global eth0
       valid_lft forever preferred_lft forever
[root@master1 3.2]# 

```


## 验证http容器功能


```
[root@master1 3.2]# curl 172.17.0.2
********** go version **********
go1.17.6
********** request header **********
User-Agent=[curl/7.29.0]
Accept=[*/*]
********** INFO  **********

```

## 上传镜像
```
[root@master1 3.2]# docker image tag go-httpserver qiubinwei513/go-httpserver
[root@master1 3.2]# docker push qiubinwei513/go-httpserver
Using default tag: latest
The push refers to repository [docker.io/qiubinwei513/go-httpserver]
f7dbfd63ab09: Pushing [==================================================>]  6.157MB/6.157MB
d950a0fc2915: Pushing [==================================================>]  4.096kB
758e62195b76: Pushing  3.072kB
ff1ca36e88c3: Pushing  3.072kB
ff3ebff0c3fd: Pushing [==============>                                    ]    117MB/407.8MB
bab89c562840: Waiting 
26a504e63be4: Waiting 
8bf42db0de72: Waiting 
31892cc314cb: Waiting 
11936051f93b: Waiting 
```
