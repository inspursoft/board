FROM grafana/grafana:4.6.3

RUN grafana-cli plugins install raintank-kubernetes-app && mkdir /plugins && cp -r /var/lib/grafana/plugins/* /plugins
ENV GF_PATHS_PLUGINS=/plugins \
    GF_AUTH_ANONYMOUS_ENABLED=true \
    GF_AUTH_ANONYMOUS_ORG_ROLE=Admin \
    GF_USERS_DEFAULT_THEME=light \
    GF_ANALYTICS_REPORTING_ENABLED=false \
    GF_SERVER_ROOT_URL=/grafana

ADD make/dev/container/grafana/init.sh /init.sh
RUN chmod u+x /init.sh

RUN sed -i '/^exec/i \/init.sh' /run.sh