FROM mysql:5.6.35

WORKDIR /tmp

ADD make/release/container/db/board.sql r.sql

ADD make/release/container/db/docker-entrypoint.sh /entrypoint.sh
RUN chmod u+x /entrypoint.sh
