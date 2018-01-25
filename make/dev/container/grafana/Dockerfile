FROM grafana/grafana:4.6.3-dev

RUN grafana-cli plugins install raintank-kubernetes-app && mkdir /plugins && cp -r /var/lib/grafana/plugins/* /plugins
ENV GF_PATHS_PLUGINS=/plugins

COPY make/dev/container/grafana/config /etc/grafana/config

ADD make/dev/container/grafana/init.sh /init.sh
RUN chmod u+x /init.sh

RUN sed -i '/^exec/i \/init.sh' /run.sh