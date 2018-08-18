FROM alpine:3.5

MAINTAINER huay@inspur.com

RUN echo http://mirrors.ustc.edu.cn/alpine/v3.5/main > /etc/apk/repositories; \
echo http://mirrors.ustc.edu.cn/alpine/v3.5/community >> /etc/apk/repositories; \
apk add --no-cache openssh-client 

ADD make/release/container/apiserver/apiserver /usr/bin/apiserver
ADD make/release/container/apiserver/entrypoint.sh /usr/bin/entrypoint.sh
ADD src/apiserver/templates /usr/bin/templates
COPY VERSION /usr/bin/VERSION

WORKDIR /usr/bin/

RUN chmod u+x apiserver entrypoint.sh

ENTRYPOINT ["sh","/usr/bin/entrypoint.sh"]

VOLUME ["/usr/bin"]

EXPOSE 8088
