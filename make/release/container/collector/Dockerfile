FROM alpine:3.5

MAINTAINER huay@inspur.com

ADD make/release/container/collector/collector /usr/bin/collector

RUN chmod u+x /usr/bin/collector

WORKDIR /usr/bin

ENTRYPOINT ["/usr/bin/collector"]

EXPOSE 8080
