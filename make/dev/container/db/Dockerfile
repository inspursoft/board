FROM mysql:5.6.35

WORKDIR /tmp

ADD make/dev/container/db/board.sql r.sql

ADD make/dev/container/db/docker-entrypoint.sh /entrypoint.sh
RUN chmod u+x /entrypoint.sh
