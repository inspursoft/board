FROM grafana/grafana:6.7.2
USER root
RUN grafana-cli plugins install devopsprodigy-kubegraf-app && grafana-cli plugins install grafana-piechart-panel && mkdir /plugins && cp -r /var/lib/grafana/plugins/* /plugins

ENV GF_PATHS_PLUGINS=/plugins \
    GF_AUTH_ANONYMOUS_ENABLED=true \
    GF_AUTH_ANONYMOUS_ORG_ROLE=Admin \
    GF_ANALYTICS_REPORTING_ENABLED=false \
    GF_SERVER_ROOT_URL=/grafana \
    GF_SECURITY_ALLOW_EMBEDDING=true

ADD make/dev/container/grafana/init.sh /init.sh
RUN apk add --no-cache curl && chmod a+x /init.sh

RUN sed -i '/^exec/i \/init.sh' /run.sh
