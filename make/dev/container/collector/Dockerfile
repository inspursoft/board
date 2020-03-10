FROM golang:1.14.0

MAINTAINER wangkun_lc@inspur.com

COPY src/collector /go/src/git/inspursoft/board/src/collector
COPY src/common /go/src/git/inspursoft/board/src/common
COPY src/vendor /go/src/git/inspursoft/board/vendor
WORKDIR /go/src/git/inspursoft/board/src/collector/cmd

RUN go build -v -mod=vendor -a -o /go/bin/collector

RUN chmod u+x /go/bin/collector

WORKDIR /go/bin/
CMD ["collector"]

EXPOSE 8080
