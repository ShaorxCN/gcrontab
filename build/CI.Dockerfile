FROM golang:alpine
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \ 
    apk update
RUN apk add --no-cache git gcc musl-dev docker

ENV GOPROXY=https://goproxy.io

RUN go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

WORKDIR /project
COPY go.mod /project
COPY go.sum /project 

RUN go mod download
