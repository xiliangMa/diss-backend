#!/bin/sh

#### 初始化变量和目录
echo "=========== 0. build start........ ==========="
## 脚本所在目录
SCRIPT_FOLDER=$(cd "$(dirname "$0")";pwd)
PROJECT_DIR="$SCRIPT_FOLDER/.."
BUILD_DIR="$SCRIPT_FOLDER/../build"


#### 清除缓存
echo "=========== 1. remove cache ==========="
if [ -d $BUILD_DIR ]
then
  rm -rf "$BUILD_DIR"
fi


#### 编译二进制文件
echo "=========== 2. build  diss-backen bin ==========="
mkdir -p "$BUILD_DIR/bin"
echo "build path: $(cd $BUILD_DIR; pwd)"
cd "$PROJECT_DIR" || exit
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "$BUILD_DIR/bin/diss-backend"


#### 准备 build 文件
echo "=========== 3. cp files for docker build ==========="
cd $PROJECT_DIR
cp -r conf swagger "$BUILD_DIR"
cp entrypoint.sh script/install.sh "$BUILD_DIR"
cp docker-compose-prod.yml "$BUILD_DIR/docker-compose.yml"

#### 停止容器
echo "=========== 4. stop diss-backend and db ==========="
cd $BUILD_DIR
docker-compose down


#### 删除镜像
echo "=========== 5. remove diss-backen images ==========="
DISS_BACKEND_ID=`docker images | grep diss-backend | awk '{print $3}'`
echo $DISS_BACKEND_ID
if [[ -n "$DISS_BACKEND_ID"  ]]
then
  docker rmi  -f $DISS_BACKEND_ID
fi
NONE_IMAGES_ID=`docker images -f "dangling=true" -q`
if [[ -n "$NONE_IMAGES_ID"  ]]
then
  docker rmi $NONE_IMAGES_ID
fi

# build 镜像
echo "=========== 6. build diss-backend images ==========="
cd $PROJECT_DIR
docker-compose build --no-cache

cd $BUILD_DIR
tar -cvf diss-backend.tar.gz ./*

echo "=========== 7. remove none images ==========="
NONE_IMAGES_ID=`docker images -f "dangling=true" -q`
if [[ -n "$NONE_IMAGES_ID"  ]]
then
  docker rmi $NONE_IMAGES_ID
fi

echo "=========== 8. build success ==========="