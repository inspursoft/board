FROM golang:1.14.0

MAINTAINER liyanqing@inspur.com

RUN go version

COPY src/adminserver /go/src/github.com/inspursoft/board/src/adminserver
COPY src/common /go/src/github.com/inspursoft/board/src/common
COPY src/vendor /go/src/github.com/inspursoft/board/src/vendor
COPY VERSION /go/bin/VERSION

ENV GO111MODULE=off

WORKDIR /go/src/github.com/inspursoft/board/src/adminserver

RUN go build -v -o /go/bin/adminserver && \
    chmod u+x /go/bin/adminserver

WORKDIR /go/bin/

VOLUME ["/data/adminserver/"]

CMD ["adminserver"]

EXPOSE 8080
