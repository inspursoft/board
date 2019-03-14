FROM alpine:3.7

MAINTAINER huay@inspur.com

RUN echo http://mirrors.ustc.edu.cn/alpine/v3.7/main > /etc/apk/repositories; \
echo http://mirrors.ustc.edu.cn/alpine/v3.7/community >> /etc/apk/repositories; \
apk add --no-cache openssh openssh-client openssl

COPY helm-v2.11.0-linux-amd64.tar.gz /usr/bin/helm.tar.gz
RUN tar -zxf /usr/bin/helm.tar.gz --strip-components=1  -C /usr/bin linux-amd64/helm && \
    rm -rf /usr/bin/helm.tar.gz

ADD make/release/container/apiserver/apiserver /usr/bin/apiserver
ADD src/apiserver/templates /usr/bin/templates
COPY make/release/container/apiserver/certs/ca-certificates.crt /etc/ssl/certs
COPY VERSION /usr/bin/VERSION

WORKDIR /usr/bin/

VOLUME ["/usr/bin", "/repos", "/keys"]

CMD ["apiserver"]

EXPOSE 8088
