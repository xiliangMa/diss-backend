## run
```
    1. set env:
        MARIADB_DATABASE=diss
        MARIADB_USER=diss
        MARIADB_PASSWORD=pwd
        MARIADB_HOST=DBIP
    
    2. 编译:
        go build -o bin/diss-backend
        
    3. 启动
        3.1 编译文件启动 ./diss-backend    或者  bee run -gendoc=true -downdoc=true
        
```
    

## build image
    
    go build -o bin/diss-backend
 
    macos 系统:
        GOOS=linux GOARCH=amd64 go build -o bin/diss-backend
            
    build 镜像:
        docker build -t diss-backend-1.0.0 .
    
    启动:
        docker run --name=diss-backend -d -e MARIADB_DATABASE=diss -e MARIADB_USER=diss -e MARIADB_PASSWORD=abc123 -e MARIADB_HOST=122.51.240.195  -p 8080:8080 -p 10443:10443 -p 8889:8889 diss-backend-1.0.0



## run in kubernetes cluster

    kubectl apply -f deploy-k8s.yml
