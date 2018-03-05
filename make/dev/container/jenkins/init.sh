#!/bin/bash

function util_done()
{
    local sleeptime=$1
    shift
    while true
    do
        "$@"
        if [ $? -eq 0 ]; then
            echo "xxxx"
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
        util_done $sleeptime curl http://10.164.17.34:8085/job/base
}

init &
