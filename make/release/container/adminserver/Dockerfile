FROM alpine:3.7

RUN echo http://mirrors.ustc.edu.cn/alpine/v3.7/main > /etc/apk/repositories; \
echo http://mirrors.ustc.edu.cn/alpine/v3.7/community >> /etc/apk/repositories; \
apk add --no-cache openssh openssh-client openssl docker

ADD make/release/container/adminserver/adminserver /usr/bin/adminserver
COPY VERSION /usr/bin/VERSION

RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

WORKDIR /usr/bin/

CMD ["adminserver"]

EXPOSE 8080