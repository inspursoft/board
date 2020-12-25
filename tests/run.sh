#!/bin/bash
set -e
# listDeps lists packages referenced by package in $1, 
# excluding golang standard library and packages in 
# direcotry vendor
#echo $1
#source $1
local_host="`hostname --fqdn`"
local_ip=`host $local_host 2>/dev/null | awk '{print $NF}'`
export HOST_IP=$local_ip
sudo rm -rf /root/.ssh/known_hosts

function listDeps()
{
    pkg=$1
    deps=$pkg
    ds=$(echo $(go list -f '{{.Imports}}' $pkg) | sed 's/[][]//g')
    for d in $ds
    do
        if echo $d | grep -q "service/controller" && echo $d | grep -qv "vendor"
        then
            deps="$deps,$d"
        fi

    done
}
#$@pull      

# set envirnment
deps=""

dir="$( cd "$( dirname "$0"  )" && pwd  )"
echo $dir

function rungotest()
{
packages=$(go list ... | grep -v -E 'vendor|tests'|grep $1) 
echo $packages
for package in $packages
do
    listDeps $package

    echo "DEBUG: testing package $package"
    echo "$deps"
    
    echo "---------------------------------------"
    echo $deps
    echo "+++++++++++++++++++++++++++++++++++++++"
    
    #go env used docker container
    echo "go test -v -race -coverprofile=profile.out -covermode=atomic -coverpkg $deps $package"
    go test -v -race -coverprofile=profile.out -covermode=atomic -coverpkg "$deps" $package
    if [ -f profile.out ]
    then
        cat profile.out >> coverage.txt
        rm profile.out
    fi

done
}

rungotest src/apiserver 
rungotest tokenserver 
