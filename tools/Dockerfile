FROM golang:1.14.1-alpine

MAINTAINER liyanqing@inspur.com

RUN apk update \
    && apk add --no-cache build-base\
    && rm -rf /var/cache/apk/*

ENV GO111MODULE=off

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone