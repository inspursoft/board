FROM golang:1.9.6

MAINTAINER wangkun_lc@inspur.com

COPY src/apiserver /go/src/git/inspursoft/board/src/apiserver
COPY src/common /go/src/git/inspursoft/board/src/common
COPY src/vendor /go/src/git/inspursoft/board/vendor
COPY src/apiserver/templates /go/bin/templates
COPY VERSION /go/bin/VERSION

WORKDIR /go/src/git/inspursoft/board/src/apiserver

RUN go build -v -a -o /go/bin/apiserver && \
    chmod u+x /go/bin/apiserver && \
    mkdir repos && \
    mkdir keys

WORKDIR /go/bin/

VOLUME ["/go/bin", "/repos", "/keys"]

CMD ["apiserver"]

EXPOSE 8088
