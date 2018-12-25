FROM library/ubuntu:14.04

RUN update-locale && rm /etc/rsyslog.d/* && rm /etc/rsyslog.conf

ADD make/dev/container/log/rsyslog.conf /etc/rsyslog.conf

# rotate logs weekly
# notes: file name cannot contain dot, or the script will not run
ADD make/dev/container/log/rotate.sh /etc/cron.weekly/rotate

# rsyslog configuration file for docker
ADD make/dev/container/log/rsyslog_docker.conf /etc/rsyslog.d/

VOLUME /var/log/board/

EXPOSE 514

CMD cron && rsyslogd -n
