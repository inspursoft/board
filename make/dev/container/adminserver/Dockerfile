FROM golang:1.13.5-alpine3.10

COPY src/adminserver /go/src/git/inspursoft/board/src/adminserver
 
WORKDIR /go/src/git/inspursoft/board/src/adminserver/

RUN apk add git \
    && go get github.com/astaxie/beego \
    && go get github.com/beego/bee
    
VOLUME ["/go/cfgfile"]

CMD [ "bee", "run", "-downdoc=true", "-gendoc=true" ]

EXPOSE 8080
