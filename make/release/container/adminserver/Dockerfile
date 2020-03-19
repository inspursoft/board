FROM golang:alpine

COPY src/adminserver /go/src/git/inspursoft/board/src/adminserver
COPY src/vendor /go/src/git/inspursoft/board/src/vendor
COPY VERSION /go/bin/VERSION

RUN apk update \
    && apk add --no-cache build-base\
    && rm -rf /var/cache/apk/*

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct

WORKDIR /go/src/git/inspursoft/board/src/

RUN go mod init && \
    go mod vendor

WORKDIR /go/src/git/inspursoft/board/src/adminserver

RUN go build -mod=vendor -v -o /go/bin/adminserver && \
    chmod u+x /go/bin/adminserver

VOLUME ["/go/cfgfile"]

WORKDIR /go/bin/

CMD ["adminserver"]

EXPOSE 8080
