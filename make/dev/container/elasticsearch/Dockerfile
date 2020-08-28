FROM docker.elastic.co/elasticsearch/elasticsearch:7.7.0

ENV CERT_DIR=/usr/share/elasticsearch/config/certs
ENV TAKE_FILE_OWNERSHIP=true \
    bootstrap.memory_lock=true \
    node.name=board-es \
    cluster.name=board-elasticsearch-cluster \
    cluster.initial_master_nodes=board-es \
    node.ml=false \
    xpack.ml.enabled=false \
    xpack.monitoring.collection.enabled=false \
    xpack.security.enabled=true \
    xpack.security.transport.ssl.enabled=true \
    xpack.security.transport.ssl.verification_mode=certificate \
    xpack.security.transport.ssl.certificate_authorities=$CERT_DIR/ca/ca.crt \
    xpack.security.transport.ssl.certificate=$CERT_DIR/elasticsearch/elasticsearch.crt \
    xpack.security.transport.ssl.key=$CERT_DIR/elasticsearch/elasticsearch.key

#USER root
COPY make/dev/container/elasticsearch/docker-entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]