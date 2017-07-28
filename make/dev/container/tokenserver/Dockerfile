FROM golang:1.8.1

MAINTAINER huay@inspur.com

COPY src/. /go/src/git/inspursoft/board/src/.

WORKDIR /go/src/git/inspursoft/board/src/tokenserver
COPY src/tokenserver/app.conf /go/bin/app.conf

RUN go build -v -a -o /go/bin/tokenserver

RUN chmod u+x /go/bin/tokenserver

WORKDIR /go/bin/
ENTRYPOINT ["/go/bin/tokenserver"]

EXPOSE 4000
