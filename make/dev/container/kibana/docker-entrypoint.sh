#!/bin/bash

ELASTICSEARCH_URL=${ELASTICSEARCH_HOSTS/[\"/}
ELASTICSEARCH_URL=${ELASTICSEARCH_URL/\"]/}
if [[ "${ELASTICSEARCH_USERNAME}" != "" && "${ELASTICSEARCH_PASSWORD}" != "" ]]; then
    ELASTICSEARCH_URL="-u ${ELASTICSEARCH_USERNAME}:${ELASTICSEARCH_PASSWORD} ${ELASTICSEARCH_URL}"
fi

echo "ELASTICSEARCH_URL=${ELASTICSEARCH_URL}"
# max wait 5min for elasticsearch start, then it will be quit
i=1
while true;
do
    status=$(curl ${ELASTICSEARCH_URL}/_cat/health?pretty 2>/dev/null | awk '{print $4}')
    if [ "$status" == "green" ] || [ "$status" == "yellow" ]; then
        break
    fi
    if [ $i -gt 300 ]; then
        echo "The kibana exit, make sure the elasticsearch is ok"
        exit 1
    fi
    echo "Waiting for elasticsearch..."
    sleep 1
    let i++
done

set -xe

# save the kibana metadata if it does not exist
found=$(curl ${ELASTICSEARCH_URL}/board/_doc/kibana?pretty 2>/dev/null | awk '/found/ {print $3}' | tr -d ,)
if [ "$found" != "true" ]; then
    # create .kibana_1 index
    code=$(curl ${ELASTICSEARCH_URL}/.kibana_1?pretty 2>/dev/null | awk '/status/ {print $3}')
    if [ "$code" == "404" ]; then
        curl -XPUT -H 'Content-Type: application/json' ${ELASTICSEARCH_URL}/.kibana_1?pretty -d @/config/index/kibana.json 2>/dev/null
    fi
    # add logstash-* index pattern
    curl -XPUT -H 'Content-Type: application/json' ${ELASTICSEARCH_URL}/.kibana_1/_doc/index-pattern:board-index?pretty -d @/config/index/pattern.json 2>/dev/null
    # set default log pattern to logstash-*
    curl -XPUT -H 'Content-Type: application/json' ${ELASTICSEARCH_URL}/.kibana_1/_doc/config:7.7.0?pretty -d @/config/index/config.json 2>/dev/null

    ###### dashboard 1: log dashboard
    # discover
    curl -XPUT -H 'Content-Type: application/json' ${ELASTICSEARCH_URL}/.kibana_1/_doc/search:logs?pretty -d @/config/discover/logs.json 2>/dev/null
    # visualization
    curl -XPUT -H 'Content-Type: application/json' ${ELASTICSEARCH_URL}/.kibana_1/_doc/visualization:logs-gauge?pretty -d @/config/visualize/logs-gauge.json 2>/dev/null
    curl -XPUT -H 'Content-Type: application/json' ${ELASTICSEARCH_URL}/.kibana_1/_doc/visualization:logs-area?pretty -d @/config/visualize/logs-area.json 2>/dev/null
    curl -XPUT -H 'Content-Type: application/json' ${ELASTICSEARCH_URL}/.kibana_1/_doc/visualization:logs-line?pretty -d @/config/visualize/logs-line.json 2>/dev/null
    # dashboard
    curl -XPUT -H 'Content-Type: application/json' ${ELASTICSEARCH_URL}/.kibana_1/_doc/dashboard:podlogs?pretty -d @/config/dashboard/logs.json 2>/dev/null
    ######

    ###### dashboard 2: error dashboard
    # discover
    curl -XPUT -H 'Content-Type: application/json' ${ELASTICSEARCH_URL}/.kibana_1/_doc/search:errors?pretty -d @/config/discover/errors.json 2>/dev/null
    # visualization
    curl -XPUT -H 'Content-Type: application/json' ${ELASTICSEARCH_URL}/.kibana_1/_doc/visualization:errors-metric?pretty -d @/config/visualize/errors-metric.json 2>/dev/null
    curl -XPUT -H 'Content-Type: application/json' ${ELASTICSEARCH_URL}/.kibana_1/_doc/visualization:errors-pie?pretty -d @/config/visualize/errors-pie.json 2>/dev/null
    curl -XPUT -H 'Content-Type: application/json' ${ELASTICSEARCH_URL}/.kibana_1/_doc/visualization:errors-line?pretty -d @/config/visualize/errors-line.json 2>/dev/null
    curl -XPUT -H 'Content-Type: application/json' ${ELASTICSEARCH_URL}/.kibana_1/_doc/visualization:errors-line-bar?pretty -d @/config/visualize/errors-line-bar.json 2>/dev/null
    # dashboard
    curl -XPUT -H 'Content-Type: application/json' ${ELASTICSEARCH_URL}/.kibana_1/_doc/dashboard:poderrors?pretty -d @/config/dashboard/errors.json 2>/dev/null
    ######

    ###### dashboard 3: log_zh_CN dashboard
    # dashboard
    curl -XPUT -H 'Content-Type: application/json' ${ELASTICSEARCH_URL}/.kibana_1/_doc/dashboard:podlogs_zh_CN?pretty -d @/config/dashboard/logs_zh_CN.json 2>/dev/null
    ######

    ###### dashboard 4: error_zh_CN dashboard
    curl -XPUT -H 'Content-Type: application/json' ${ELASTICSEARCH_URL}/.kibana_1/_doc/dashboard:poderrors_zh_CN?pretty -d @/config/dashboard/errors_zh_CN.json 2>/dev/null
    ######


    # kibana metadata saves sucessfully tag
    curl -XPUT -H 'Content-Type: application/json' ${ELASTICSEARCH_URL}/board/_doc/kibana?pretty -d '{"successfully": true}' 2>/dev/null
fi

exec su kibana -c "/bin/bash /usr/local/bin/kibana-docker"