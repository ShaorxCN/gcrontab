#!/usr/bin/env bash -e

#运行所有单元测试
go test ../...
read -n1 -p "单元测试完成，是否要继续 [Y/N]? " answer
case $answer in
    Y | y)
        echo -e "\n开始编译";;
    *)
        exit 0;;
esac

#编译linux执行文件,交叉编译不支持cgo所以禁用
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gcrontab ../main.go

echo "编译完成"

#编译docker镜像
docker build -t gcrontab:dev-latest -f ./Dockerfile .
