FROM golang:1.14.0

MAINTAINER huay@inspur.com

COPY src/. /go/src/github.com/inspursoft/board/src/.

WORKDIR /go/src/github.com/inspursoft/board/src/tokenserver

ENV GO111MODULE=off

RUN go build -v -a -o /go/bin/tokenserver

RUN chmod u+x /go/bin/tokenserver

WORKDIR /go/bin/

ENTRYPOINT ["/go/bin/tokenserver"]

VOLUME ["/go/bin"]

EXPOSE 4000
