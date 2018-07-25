FROM alpine:3.7

MAINTAINER huay@inspur.com

RUN echo http://mirrors.ustc.edu.cn/alpine/v3.7/main > /etc/apk/repositories; \
echo http://mirrors.ustc.edu.cn/alpine/v3.7/community >> /etc/apk/repositories; \
apk add --no-cache openssh openssh-client 

ADD make/release/container/apiserver/apiserver /usr/bin/apiserver
ADD src/apiserver/templates /usr/bin/templates
COPY make/release/container/apiserver/certs/ca-certificates.crt /etc/ssl/certs
COPY VERSION /usr/bin/VERSION

WORKDIR /usr/bin/

VOLUME ["/usr/bin", "/repos", "/keys"]

CMD ["apiserver"]

EXPOSE 8088
