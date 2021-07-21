#!/bin/sh

ARCH=$1
if [ ! -n "$ARCH" ]
then
  ARCH=amd64
fi

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
echo "=========== 2. build  diss-backend bin ==========="
mkdir -p "$BUILD_DIR/bin"
echo "build path: $(cd $BUILD_DIR; pwd)"
cd "$PROJECT_DIR" || exit
sed -is 's/RunMode = dev/RunMode = prod/' "$PROJECT_DIR/conf/app.conf" && rm "$PROJECT_DIR/conf/app.confs"
if [ $ARCH == 'arm' ]
then
  CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o "$BUILD_DIR/bin/diss-backend"
else
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "$BUILD_DIR/bin/diss-backend"
fi


#### 准备 build 文件
echo "=========== 3. cp files for docker build ==========="
cd $PROJECT_DIR
cp -r conf swagger "$BUILD_DIR"
rm -rf "$BUILD_DIR/conf/app.conf"
mv "$BUILD_DIR/conf/app-prod.conf" "$BUILD_DIR/conf/app.conf"
cp entrypoint.sh "$BUILD_DIR"
if [ $ARCH == 'arm' ]
then
  cp docker-compose-arm.yml "$BUILD_DIR/docker-compose.yml"
else
  cp docker-compose.yml "$BUILD_DIR/docker-compose.yml"
fi
mkdir -p "$BUILD_DIR/upload/plugin/scope"
mkdir -p "$BUILD_DIR/upload/license"
mkdir -p "$BUILD_DIR/upload/logo"
mkdir -p "$BUILD_DIR/public"
cp ./upload/plugin/scope/diss-scope.yml "$BUILD_DIR/upload/plugin/scope"
cp ./upload/license/TrialLicense.lic "$BUILD_DIR/upload/license/TrialLicense.lic"
cp ./upload/logo/newcon.png "$BUILD_DIR/upload/logo/newcon.png"
cp ./public/apm "$BUILD_DIR/public/apm"


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
if [ $ARCH == 'arm' ]
then
  docker-compose -f docker-compose-arm.yml build --no-cache
else
  docker-compose build --no-cache
fi

cd $BUILD_DIR
tar -zcvf diss-backend.tar.gz ./docker-compose.yml ./conf
rm -rf bin swagger entrypoint.sh upload public

echo "=========== 7. remove none images ==========="
NONE_IMAGES_ID=`docker images -f "dangling=true" -q`
if [[ -n "$NONE_IMAGES_ID"  ]]
then
  docker rmi $NONE_IMAGES_ID
fi

echo "=========== 8. build success ==========="
sed -is 's/RunMode = prod/RunMode = dev/' "$PROJECT_DIR/conf/app.conf" && rm "$PROJECT_DIR/conf/app.confs"