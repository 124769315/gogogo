f**MOdule8作业**
# 一、配置Yaml文件
## vim http.yaml
```
root@master01:/m8# cat http.yaml 
apiVersion: v1
kind: Pod
metadata:
        name: http-qbw03
spec:
        containers:
                - name: http-qbw03
                  image: localhost:5000/http
                  imagePullPolicy: IfNotPresent
                  resources:
                          limits:
                                  cpu: 700m
                                  memory: "500Mi"
                          requests:
                                  cpu: 700m
                                  memory: "500Mi"
.
                  readinessProbe:
                          httpGet:        
                                  path: /healthz
                                  port: 80
                          initialDelaySeconds: 15
                          periodSeconds: 5
                  lifecycle:
                          preStop:
                                  exec:
                                          command: ["PID=`ps -ef | grep httpser | grep -v sh | grep -v grep | awk '{print $2}'` && kill -9 $PID"]
```
### 1 YAML参数备注
#### 1.1  配置本地镜像
##### imagePullPolicy: IfNotPresent
```
优先读取该节点本地image镜像
```
##### 上传本地http image至docker仓库

##### 配置k8s本地image方法:
```
1. 搭建Docker Registry
This image contains an implementation of the Docker Registry HTTP API V2 for use with Docker 1.6+. See github.com/docker/distribution for more details about what it is.

Run a local registry: Quick Version

$ docker run -d -p 5000:5000 --restart always --name registry registry:2
Now, use it from within Docker:

$ docker tag account-svc localhost:5000/account-svc
$ docker push localhost:5000/account-svc
查看镜像:

curl http://localhost:5000/v2/_catalog

确认有http镜像
root@master01:~# curl http://localhost:5000/v2/_catalog
{"repositories":["http"]}

```
#### 设置qos
##### 当limits=requests时，pods的qosClass为Guaranteed
```
resources:

    limits:
    
        cpu: 700m
        memory: "500Mi"
        
    requests:
    
        cpu: 700m
        memory: "500Mi"
```
##### readinessProbe 优雅启动pod+探活，通过httpGet实现
```
                  readinessProbe:
                          httpGet:        
                                  path: /healthz
                                  port: 80
                          initialDelaySeconds: 15
                          periodSeconds: 5
```

##### preStop 优化停止pod，通过执行kill pid实现
```
                  lifecycle:
                          preStop:
                                  exec:
                                          command: ["PID=`ps -ef | grep httpser | grep -v sh | grep -v grep | awk '{print $2}'` && kill -9 $PID"]
```
# 二、运行验证
## 启动pod
### kubectl create -f http.yaml
```
root@master01:/m8# kubectl create -f http.yaml 
pod/http-qbw03 created
root@master01:/m8#
```
### watch查看container在1秒内完成创建，处于Ready状态为0/1，pod起了但是由于优雅启动机制在，因此没有立刻对外提供服务。延迟启动15秒后，探活机制生效，过了5秒httpGet获取到了200请求，Ready状态变为1/1， 对外提供服务
```
root@master01:~# kubectl get pod -w
NAME                               READY   STATUS    RESTARTS      AGE
envoy-6958c489d9-7jf79             1/1     Running   1 (26h ago)   3d4h
ng-deployment                      1/1     Running   2 (26h ago)   25d
nginx-deployment-8d545c96d-ls7sm   1/1     Running   1 (26h ago)   3d4h
http-qbw03                         0/1     Pending   0             0s
http-qbw03                         0/1     Pending   0             0s
http-qbw03                         0/1     ContainerCreating   0             0s
http-qbw03                         0/1     ContainerCreating   0             1s
http-qbw03                         0/1     Running             0             1s
http-qbw03                         1/1     Running             0             20s

```


# 创建deployment
![image.png](https://note.youdao.com/yws/res/2/WEBRESOURCE09e067007e05383dad9fb443bfc43692)
```
apiVersion: apps/v1
kind: Deployment
metadata:
        name: http-deployment
spec:
        replicas: 3
        selector:
                matchLabels:
                        app: nginx
        template:
                        metadata:
                                labels:
                                        app: nginx
                        spec:
                                containers:
                                      - name: http-deployment-qbw
                                        image: localhost:5000/http
                                        imagePullPolicy: IfNotPresent
                                        resources:
                                                limits:
                                                        cpu: 100m
                                                        memory: "100Mi"
                                                requests:
                                                        cpu: 100m
                                                        memory: "100Mi"
                                        readinessProbe:
                                                httpGet:
                                                        path: /healthz
                                                        port: 80
                                                initialDelaySeconds: 15
                                                periodSeconds: 5
                                        lifecycle:
                                                preStop:
                                                        exec:
                                                                command: ["PID=`ps -ef | grep httpser | grep -v sh | grep -v grep | awk '{print $2}'` && kill -9 $PID"]

```

## 运行deployment
![image.png](https://note.youdao.com/yws/res/0/WEBRESOURCE204a9b810d2070ebbf4cc5b248a0e550)

## 创建对应service
![image.png](https://note.youdao.com/yws/res/e/WEBRESOURCEd89cb6d36b65ab83ad2e08a3b6e687ee)
![image.png](https://note.youdao.com/yws/res/7/WEBRESOURCE3c306213088ba86e4b45941a39f279f7)
## 创建ingress
![image.png](https://note.youdao.com/yws/res/d/WEBRESOURCE93deef5f295f95760ea8d2c8ee5e36ad)
![image.png](https://note.youdao.com/yws/res/a/WEBRESOURCE7d581f5a46807b040fa27c2b51c722ca)

## 通过https + curl 访问
![image.png](https://note.youdao.com/yws/res/4/WEBRESOURCE8be328b9d930633316871e861e202444)



# 创建deployment
```
apiVersion: apps/v1
kind: Deployment
metadata:
        name: http-deployment
spec:
        replicas: 3
        selector:
                matchLabels:
                        app: nginx
        template:
                        metadata:
                                labels:
                                        app: nginx
                        spec:
                                containers:
                                      - name: http-deployment-qbw
                                        image: localhost:5000/http
                                        imagePullPolicy: IfNotPresent
                                        resources:
                                                limits:
                                                        cpu: 100m
                                                        memory: "100Mi"
                                                requests:
                                                        cpu: 100m
                                                        memory: "100Mi"
                                        readinessProbe:
                                                httpGet:
                                                        path: /healthz
                                                        port: 80
                                                initialDelaySeconds: 15
                                                periodSeconds: 5
                                        lifecycle:
                                                preStop:
                                                        exec:
                                                                command: ["PID=`ps -ef | grep httpser | grep -v sh | grep -v grep | awk '{print $2}'` && kill -9 $PID"]

```

## 运行deployment
```
root@master01:/m8# k get pod
NAME                              READY   STATUS    RESTARTS      AGE
envoy-6958c489d9-7jf79            1/1     Running   1 (11d ago)   13d
http-deployment-c859ff4d7-fbdrb   1/1     Running   0             31m
http-deployment-c859ff4d7-gpxp5   1/1     Running   0             8m5s
http-deployment-c859ff4d7-rnbs7   1/1     Running   0             31m
jenkins-0                         1/1     Running   0             40h
ng-deployment                     1/1     Running   2 (11d ago)   35d
```

## 创建对应service
```
root@master01:/m8# cat service.yaml 
apiVersion: v1
kind: Service
metadata:
        name: http-basic
spec:
        type: ClusterIP
        ports:
                - port: 80
                  protocol: TCP
                  name: http
        selector:
                app: nginx

```
### 查看svc
```
root@master01:/m8# k get svc
NAME             TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)                        AGE
envoy            NodePort    10.102.222.191   <none>        10000:31565/TCP                42d
http-basic       ClusterIP   10.111.26.215    <none>        80/TCP                         8d
jenkins          NodePort    10.104.153.112   <none>        80:31364/TCP,50000:32456/TCP   43h
kubernetes       ClusterIP   10.96.0.1        <none>        443/TCP                        43d
nginx-service    NodePort    10.109.1.139     <none>        80:30080/TCP                   43d
readiness-gate   ClusterIP   10.110.144.224   <none>        80/TCP                         13d

```
## 安装ingress

## 创建ingress
```
root@master01:/m8# cat http-ingress.yaml 
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: gateway
  annotations:
    kubernetes.io/ingress.class: "nginx"
spec:
  tls:
    - hosts:
        - cncamp.com
      secretName: cncamp-tls
  rules:
    - host: cncamp.com
      http:
        paths:
          - path: "/"
            pathType: Prefix
            backend:
              service:
                name: http-basic
                port:
                  number: 80
```
### 查看ingress
```
root@master01:/m8# k get ingress
NAME      CLASS    HOSTS        ADDRESS     PORTS     AGE
gateway   <none>   cncamp.com   10.0.16.8   80, 443   7d9h
```
### 查看ingress对应外部controller IP及端口
```
root@master01:/m8# kubectl get svc -n ingress-nginx
NAME                                 TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)                      AGE
ingress-nginx-controller             NodePort    10.99.39.141   <none>        80:31233/TCP,443:30054/TCP   7d12h
ingress-nginx-controller-admission   ClusterIP   10.103.78.85   <none>        443/TCP                      7d12h

```

## 通过https + curl 访问
```
root@master01:/m8# curl -H "Host: cncamp.com" https://101.34.27.159:30054  -k
********** go version **********
go1.13.8
********** request header **********
X-Real-Ip=[172.16.241.64]
X-Forwarded-For=[172.16.241.64]
X-Forwarded-Host=[cncamp.com]
X-Forwarded-Proto=[https]
X-Forwarded-Scheme=[https]
User-Agent=[curl/7.68.0]
Accept=[*/*]
X-Request-Id=[087cb3e9ee06bac7e0f104a9da5147d5]
X-Forwarded-Port=[443]
X-Scheme=[https]
********** INFO  **********

```

