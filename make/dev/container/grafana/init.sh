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

function init()
{
    if [ ! -f /var/lib/grafana/initok ]; then
        #replace environment variables
        if [ ! -d /etc/grafana/config/ ]; then
            echo "directory /etc/grafana/config/ does not exist. make sure you have executed the 'make prepare' command"
            exit 1
        fi

        allfiles=`ls -l /etc/grafana/config/ | wc -l`
        if [ $allfiles == "1" ]; then
            echo "directory /etc/grafana/config/ does not have files. make sure you have executed the 'make prepare' command"
            exit 2
        fi

        sleeptime=5
        # check server status
        echo "checking grafana server"
        util_done $sleeptime curl ${GRAFANA_ADDRESS}/api/org
        echo "grafana server has already started"

        # set the kubernetes plugin
        echo "setting the kubernetes plugin"
        util_done $sleeptime curl -X POST -H "Content-Type: application/json;charset=UTF-8" -d @/etc/grafana/config/kubernetes.json ${GRAFANA_ADDRESS}/api/plugins/raintank-kubernetes-app/settings
        echo "set the kubernetes plugin successfully"

        # add the graphite datasource
        echo "adding the graphite datasource to grafana"
        util_done $sleeptime curl -X POST -H "Content-Type: application/json;charset=UTF-8" -d @/etc/grafana/config/graphite-datasource.json ${GRAFANA_ADDRESS}/api/datasources
        echo "add the graphite datasource to grafana successfully"

        # add the kubernetes datasource
        echo "adding the kubernetes datasource to grafana"
        util_done $sleeptime curl -X POST -H "Content-Type: application/json;charset=UTF-8" -d @/etc/grafana/config/kubernetes-datasource.json ${GRAFANA_ADDRESS}/api/datasources
        echo "add the kubernetes datasource to grafana successfully"

        # add the kuberenets dashboard
        echo "adding the kubernetes dashboard to grafana"
        util_done $sleeptime curl -X POST -H "Content-Type: application/json;charset=UTF-8" -d @/etc/grafana/config/kubernetes-dashboard.json ${GRAFANA_ADDRESS}/api/dashboards/db
        echo "add the kubernetes dashboard to grafana successfully"

        # generate the install tag file 
        echo "init successfully"
        echo "init successfully" > /var/lib/grafana/initok
    fi
}

init &