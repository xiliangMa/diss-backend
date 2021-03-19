# diss-backend
diss service

## dev
```.shell script
bee run -gendoc=true -downdoc=true
```   

## build image
```
./script/build.sh
```

## run
启动
```
cd build 
docker-compose up -d
```

关闭：
```
docker-compose down -v
```
