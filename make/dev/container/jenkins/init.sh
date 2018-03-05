#!/bin/bash

function util_done()
{
    local sleeptime=$1
    shift
    while true
    do
        "$@"
        if [ $? -eq 0 ]; then
            echo "start to add node ......"
             /usr/share/jenkins/addnode.sh
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
