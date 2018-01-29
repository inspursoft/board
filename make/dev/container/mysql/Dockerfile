FROM mysql:5.6.35

WORKDIR /tmp

ADD make/dev/container/mysql/board.sql r.sql

ADD make/dev/container/mysql/docker-entrypoint.sh /entrypoint.sh
RUN chmod u+x /entrypoint.sh
