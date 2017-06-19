FROM golang:1.8.1

MAINTAINER wangkun_lc@inspur.com

COPY src/apiserver /go/src/apiserver
WORKDIR /go/src/apiserver

RUN go build -v -a -o /go/bin/apiserver

RUN chmod u+x /go/bin/apiserver

WORKDIR /go/bin/
CMD ["apiserver"]

EXPOSE 8080
