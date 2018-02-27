FROM grafana/grafana:4.6.3

RUN grafana-cli plugins install raintank-kubernetes-app && mkdir /plugins && cp -r /var/lib/grafana/plugins/* /plugins
ENV GF_PATHS_PLUGINS=/plugins

ADD make/dev/container/grafana/init.sh /init.sh
RUN chmod u+x /init.sh

RUN sed -i '/^exec/i \/init.sh' /run.sh