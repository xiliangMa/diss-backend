#!/bin/sh
echo "=========== 0. build start........ ==========="

# 编译二进制文件
echo "=========== 1. buid  diss-backen ==========="
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/diss-backend

#清除缓存
echo "=========== 2. remove cache ==========="
rm -rf build/*


# 准备 build 文件
echo "=========== 3. cp build file ==========="
cp -r conf swagger bin build
cp entrypoint.sh bin/diss-backend script/install.sh build
cp docker-compose-prod.yml build/docker-compose.yml

#停止容器
echo "=========== 4. stop diss-backend and db ==========="
docker-compose down

#删除镜像
echo "=========== 5. remove diss-backen images ==========="
docker rmi `docker images | grep diss-backend | awk '{print $3}'` -f
docker rmi `docker images -f "dangling=true" -q`

# build 镜像
echo "=========== 6. build diss-backend images ==========="
docker-compose build --no-cache

cd build
tar -cvf diss-backend.tar.gz ./*
echo "=========== 7. build success ==========="