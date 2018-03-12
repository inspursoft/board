#!/bin/bash

function util_done()
{
    local sleeptime=$1
    shift
    while true
    do
        "$@"
        if [ $? -eq 0 ]; then
            echo "add node to $jenkins_node_ip"
             python /usr/share/jenkins/addNode.py
            break
        fi
        sleep $sleeptime
    done
}

function init()
{
        sleeptime=5
        # check server status
        echo "checking jenkins server"
        util_done $sleeptime curl http://$jenkins_host_ip:$jenkins_host_port/job/base
}

init &
