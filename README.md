# diss-backend
diss-backend

## 开发环境

set env:
```
MARIADB_DATABASE=<DB>
MARIADB_USER=<DB_User>
MARIADB_PASSWORD=<DB_PWD>
MARIADB_HOST=<DB_IP>
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
-e MARIADB_DATABASE=<DB> \
-e MARIADB_USER=<DB_USER> \
-e MARIADB_PASSWORD=<DB_PWD> \
-e MARIADB_HOST=<DB_IP>  \
-p 8080:8080 \
-p 10443:10443 \
-p 8889:8889 \
diss-backend-1.0.0
```


## Run in kubernetes cluster

kubectl apply -f deploy-k8s.yml
