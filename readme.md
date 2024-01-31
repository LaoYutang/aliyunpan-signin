# 阿里云盘自动签到

## 原理
使用阿里云盘refresh-token获取access-token，调用接口签到并且获取奖励。定时时间为每天2点钟。

## docker启动
> dockerhub镜像只支持amd64架构，其他架构请自行打包镜像
```bash
docker run -d \
  --name alisignin \
  -v /etc/localtime:/etc/localtime \
  -e refresh-token=XXXXXXXXX \
  laoyutang/aliyunpan-signin:latest
```

## 打包
```bash
docker build .
```

## go交叉编译
```bash
GOOS=linux CGO_ENABLED=0 GOARCH=arm64 go build .
```