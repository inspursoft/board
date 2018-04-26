FROM docker.elastic.co/kibana/kibana-oss:6.2.4

ENV SERVER_BASEPATH=/kibana \
    ELASTICSEARCH_URL=http://elasticsearch:9200

USER root
COPY make/release/container/kibana/docker-entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

CMD ["/entrypoint.sh"]