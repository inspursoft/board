FROM golang:1.8.1

MAINTAINER wangkun_lc@inspur.com

COPY src/apiserver /go/src/git/inspursoft/board/src/apiserver
COPY src/common /go/src/git/inspursoft/board/src/common
COPY src/vendor /go/src/git/inspursoft/board/vendor

WORKDIR /go/src/git/inspursoft/board/src/apiserver

RUN go build -v -a -o /go/bin/apiserver

RUN chmod u+x /go/bin/apiserver

WORKDIR /go/bin/

CMD ["apiserver"]

VOLUME ["/go/bin"]

EXPOSE 8080
