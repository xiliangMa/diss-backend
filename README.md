# diss-backend
diss-backend

## 开发环境
> 前提：安装 mysql 5.7+
 
set env:
```
MYSQL_DATABASE=<DB>
MYSQL_USER=<DB_User>
MYSQL_PASSWORD=<DB_PWD>
MYSQL_HOST=<DB_IP>
```
    
build:
> CGO_ENABLED=0 (alpine容器无法运行go编译的二进制文件问题解决)
```
CGO_ENABLED=0 go build -o bin/diss-backend
```

run:

```
编译文件启动 ./diss-backend    或者  bee run -gendoc=true -downdoc=true
```
    
## Run in Docker

build:
 > CGO_ENABLED=0 (alpine容器无法运行go编译的二进制文件问题解决)
```
CGO_ENABLED=0 go build -o bin/diss-backend

macos 系统:
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/diss-backend
```
 
            
build image:
 ```
docker build -t diss-backend .
```
> 前提：安装mysql 5.7+

修改 diss-backend.env:
```
MYSQL_ROOT_PASSWORD=diss
MYSQL_USER=diss
MYSQL_PASSWORD=diss
MYSQL_DATABASE=diss
MYSQL_HOST=diss-db
```


run:
```
docker run --name=diss-backend -d \
--env-file  ./deploy/diss-backend.env \
--env-file  ./deploy/diss-backend-db.env \
-p 8080:8080 \
-p 10443:10443 \
-p 8889:8889 \
diss-backend:latest
```



## Run by docker-compose
启动
```
docker-compose up -d
```

关闭：
```
docker-compose down
```


## Run by kubernetes
```
cd deploy/kubernetes
kubectl apply -f install.yml
```