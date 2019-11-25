# diss-backend
diss-backend

## 开发环境

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
docker build -t diss-backend-1.0.0 .
```
    
run:
 ```
docker run --name=diss-backend -d \
-e MYSQL_DATABASE=<DB> \
-e MYSQL_USER=<DB_USER> \
-e MYSQL_PASSWORD=<DB_PWD> \
-e MYSQL_HOST=<DB_IP>  \
-p 8080:8080 \
-p 10443:10443 \
-p 8889:8889 \
diss-backend-1.0.0
```



## Run by docker-compose
修改 diss-backend.env:
```
MYSQL_ROOT_PASSWORD=diss
MYSQL_USER=diss
MYSQL_PASSWORD=diss
MYSQL_DATABASE=diss
MYSQL_HOST=diss-db
```

启动
```
docker-compose up -d
```

关闭：
```
docker-compose down
```

## Run in kubernetes cluster

kubectl apply -f deploy-k8s.yml
