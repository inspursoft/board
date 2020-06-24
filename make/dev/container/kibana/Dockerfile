FROM docker.elastic.co/kibana/kibana:7.7.0

ENV SERVER_BASEPATH=/kibana \
    ELASTICSEARCH_HOSTS='["http://elasticsearch:9200"]' \
    XPACK_APM_ENABLED=false \
    XPACK_INFRA_ENABLED=false \
    XPACK_ML_ENABLED=false \
    MONITORING_ENABLED=false \
    XPACK_REPORTING_ENABLED=false \
    XPACK_SECURITY_ENABLED=false \
    XPACK_SPACES_ENABLED=false \
    XPACK_UPTIME_ENABLED=false \
    XPACK_SIEM_ENABLED=false \
    XPACK_MAPS_ENABLED=false \
    TELEMETRY_ENABLED=false \
    i18n.locale=zh-CN

USER root
COPY make/dev/container/kibana/docker-entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

CMD ["/entrypoint.sh"]