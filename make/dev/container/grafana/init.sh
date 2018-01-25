#!/bin/bash

function util_done()
{
    local sleeptime=$1
    shift
    while true
    do
        "$@"
        if [ $? -eq 0 ]; then
            break
        fi
        sleep $sleeptime
    done
}

function subst()
{
    perl -p -e 's/\$\{([^}]+)\}/defined $ENV{$1} ? $ENV{$1} : $&/eg; s/\$\{([^}]+)\}//eg' $1
}

function init()
{
    if [ ! -f /var/lib/grafana/initok ]; then
        #replace environment variables
        subst /etc/grafana/config/kubernetes.json 1>/tmp/kubernetes.generated
        subst /etc/grafana/config/graphite-datasource.json 1>/tmp/graphite-datasource.generated
        subst /etc/grafana/config/kubernetes-datasource.json 1>/tmp/kubernetes-datasource.generated
        subst /etc/grafana/config/kubernetes-dashboard.json 1>/tmp/kubernetes-dashboard.generated

        sleeptime=5
        # check server status
        echo "checking grafana server"
        util_done $sleeptime curl ${GRAFANA_ADDRESS}/api/org
        echo "grafana server has already started"

        # set the kubernetes plugin
        echo "setting the kubernetes plugin"
        util_done $sleeptime curl -X POST -H "Content-Type: application/json;charset=UTF-8" -d @/tmp/kubernetes.generated ${GRAFANA_ADDRESS}/api/plugins/raintank-kubernetes-app/settings
        echo "set the kubernetes plugin successfully"

        # add the graphite datasource
        echo "adding the graphite datasource to grafana"
        util_done $sleeptime curl -X POST -H "Content-Type: application/json;charset=UTF-8" -d @/tmp/graphite-datasource.generated ${GRAFANA_ADDRESS}/api/datasources
        echo "add the graphite datasource to grafana successfully"

        # add the kubernetes datasource
        echo "adding the kubernetes datasource to grafana"
        util_done $sleeptime curl -X POST -H "Content-Type: application/json;charset=UTF-8" -d @/tmp/kubernetes-datasource.generated ${GRAFANA_ADDRESS}/api/datasources
        echo "add the kubernetes datasource to grafana successfully"

        # add the kuberenets dashboard
        echo "adding the kubernetes dashboard to grafana"
        util_done $sleeptime curl -X POST -H "Content-Type: application/json;charset=UTF-8" -d @/tmp/kubernetes-dashboard.generated ${GRAFANA_ADDRESS}/api/dashboards/db
        echo "add the kubernetes dashboard to grafana successfully"

        # generate the install tag file 
        echo "init successfully"
        echo "init successfully" > /var/lib/grafana/initok
    fi
}

export GRAFANA_ADDRESS=${GRAFANA_ADDRESS-http://127.0.0.1:3000}
export GRAPHITE_PROTOCOL=${GRAPHITE_PROTOCOL-http}
export GRAPHITE_IP=${GRAPHITE_IP-graphite}
export GRAPHITE_PORT=${GRAPHITE_PORT-80}
export GRAPHITE_CARBON_PORT=${GRAPHITE_CARBON_PORT-2003}
export KUBERNETES_ADDRESS=${KUBERNESTES_ADDRESS-http://${KUBE_MASTER_IP}:${KUBE_MASTER_PORT}}

init &