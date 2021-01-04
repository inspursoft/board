FROM golang:1.14.0

MAINTAINER wangkun_lc@inspur.com

COPY src/apiserver /go/src/github.com/inspursoft/board/src/apiserver
COPY src/common /go/src/github.com/inspursoft/board/src/common
COPY src/vendor /go/src/github.com/inspursoft/board/vendor
COPY src/apiserver/templates /go/bin/templates
COPY VERSION /go/bin/VERSION
COPY helm-v2.11.0-linux-amd64.tar.gz /usr/bin/helm.tar.gz

WORKDIR /go/src/github.com/inspursoft/board/src/apiserver

ENV GO111MODULE=off

RUN go build -v -a -o /go/bin/apiserver && \
    chmod u+x /go/bin/apiserver && \
    mkdir repos && \
    mkdir keys && \
    tar -zxf /usr/bin/helm.tar.gz --strip-components=1  -C /usr/bin linux-amd64/helm && \
    rm -rf /usr/bin/helm.tar.gz

WORKDIR /go/bin/

VOLUME ["/go/bin", "/repos"]

CMD ["apiserver"]

EXPOSE 8088
