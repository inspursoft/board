FROM golang:1.9.6

COPY src/adminserver/vendor /go/src
COPY make /go/cfgfile
COPY src/adminserver /go/src/git/inspursoft/board/src/adminserver
 
WORKDIR /go/src/git/inspursoft/board/src/adminserver/

RUN go get github.com/astaxie/beego \
    && go get github.com/beego/bee \
    && rm -rf /go/src/git/inspursoft/board/src/adminserver/vendor
    
VOLUME ["/go/cfgfile", "/go/secrets"]

CMD [ "bee", "run", "-downdoc=true", "-gendoc=true" ]

EXPOSE 8080
