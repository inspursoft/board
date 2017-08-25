FROM golang:1.8.1

MAINTAINER wangkun_lc@inspur.com

COPY src/apiserver /go/src/git/inspursoft/board/src/apiserver
COPY src/common /go/src/git/inspursoft/board/src/common
COPY src/vendor /go/src/git/inspursoft/board/vendor

WORKDIR /go/src/git/inspursoft/board/src/apiserver

RUN go build -v -a -o /go/bin/apiserver && \
    chmod u+x /go/bin/apiserver

WORKDIR /go/bin/

COPY make/dev/container/apiserver/entrypoint.sh /go/bin/entrypoint.sh
RUN chmod u+x entrypoint.sh

ENTRYPOINT ["sh", "entrypoint.sh"]

VOLUME ["/go/bin"]

EXPOSE 8088
