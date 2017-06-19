FROM mysql:5.6.35

WORKDIR /tmp

ADD board.sql r.sql

ADD docker-entrypoint.sh /entrypoint.sh
RUN chmod u+x /entrypoint.sh
