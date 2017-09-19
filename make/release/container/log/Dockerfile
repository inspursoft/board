FROM alpine:3.5

RUN echo http://mirrors.ustc.edu.cn/alpine/v3.5/main > /etc/apk/repositories; \
echo http://mirrors.ustc.edu.cn/alpine/v3.5/community >> /etc/apk/repositories; \
apk add --no-cache rsyslog

#RUN runDeps='curl rsyslog' \
# && apk add --update $runDeps \
# && rm -rf /var/cache/apk/*

#RUN rm /etc/rsyslog.d/* && rm /etc/rsyslog.conf

ADD make/release/container/log/rsyslog.conf /etc/rsyslog.conf

# rotate logs weekly
# notes: file name cannot contain dot, or the script will not run
ADD make/release/container/log/rotate.sh /etc/cron.weekly/rotate

# rsyslog configuration file for docker
ADD make/release/container/log/rsyslog_docker.conf /etc/rsyslog.d/

VOLUME /var/log/board/

EXPOSE 514

CMD crond && rsyslogd -n
