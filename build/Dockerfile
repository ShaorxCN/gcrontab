# 设置继承自alpine官方镜像
FROM alpine:3.9
RUN apk add -U tzdata \
&& cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && apk add -U ca-certificates 
#添加文件到容器中
COPY ./gcrontab /

CMD  /gcrontab >> /var/log/out.log 2>&1
