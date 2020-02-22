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

```
go build -o bin/diss-backend
```

run:

```
编译文件启动 ./diss-backend    或者  bee run -gendoc=true -downdoc=true
```
    
## Run in Docker

build:
 
```
go build -o bin/diss-backend
macos 系统:
    GOOS=linux GOARCH=amd64 go build -o bin/diss-backend
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
--env-file  ./diss-backend.env \
-p 8080:8080 \
-p 10443:10443 \
-p 8889:8889 \
diss-backend
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