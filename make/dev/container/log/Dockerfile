FROM library/ubuntu:14.04

RUN rm /etc/rsyslog.d/* && rm /etc/rsyslog.conf

ADD rsyslog.conf /etc/rsyslog.conf

# rotate logs weekly
# notes: file name cannot contain dot, or the script will not run
ADD rotate.sh /etc/cron.weekly/rotate

# rsyslog configuration file for docker
ADD rsyslog_docker.conf /etc/rsyslog.d/

VOLUME /var/log/board/

EXPOSE 514

CMD cron && rsyslogd -n
