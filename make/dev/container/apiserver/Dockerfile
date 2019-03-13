FROM golang:1.9.6

MAINTAINER wangkun_lc@inspur.com

COPY src/apiserver /go/src/git/inspursoft/board/src/apiserver
COPY src/common /go/src/git/inspursoft/board/src/common
COPY src/vendor /go/src/git/inspursoft/board/vendor
COPY src/apiserver/templates /go/bin/templates
COPY VERSION /go/bin/VERSION

WORKDIR /go/src/git/inspursoft/board/src/apiserver

RUN go build -v -a -o /go/bin/apiserver && \
    chmod u+x /go/bin/apiserver && \
    mkdir repos && \
    mkdir keys && \
    curl -sLf https://storage.googleapis.com/kubernetes-helm/helm-v2.11.0-linux-amd64.tar.gz > /usr/bin/helm.tar.gz && \
    tar -zxf /usr/bin/helm.tar.gz --strip-components=1  -C /usr/bin linux-amd64/helm && \
    rm -rf /usr/bin/helm.tar.gz

WORKDIR /go/bin/

VOLUME ["/go/bin", "/repos", "/keys"]

CMD ["apiserver"]

EXPOSE 8088
