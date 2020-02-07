FROM golang:1.13.5-alpine3.10

COPY src/adminserver /go/src/board-adminserver/src
 
WORKDIR /go/src/board-adminserver/src/

RUN apk add git \
    && go get github.com/astaxie/beego \
    && go get github.com/beego/bee
    
VOLUME ["/go/cfgfile"]

CMD [ "bee", "run"]

EXPOSE 8080
