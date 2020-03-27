FROM golang:1.14.0

MAINTAINER huay@inspur.com

COPY src/. /go/src/git/inspursoft/board/src/.

WORKDIR /go/src/git/inspursoft/board/src/tokenserver

RUN go build -v -mod=vendor -a -o /go/bin/tokenserver

RUN chmod u+x /go/bin/tokenserver

WORKDIR /go/bin/

ENTRYPOINT ["/go/bin/tokenserver"]

VOLUME ["/go/bin"]

EXPOSE 4000
