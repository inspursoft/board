# listDeps lists packages referenced by package in $1, 
# excluding golang standard library and packages in 
# direcotry vendor
echo $1
source $1
local_host="`hostname --fqdn`"
local_ip=`host $local_host 2>/dev/null | awk '{print $NF}'`
export HOST_IP=$local_ip
rm -rf /root/.ssh/known_hosts

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
set -e

# set envirnment
deps=""
gopath=/go/src/git/inspursoft/board/
golangImage=golang:1.9.6-alpine3.7
volumeDir=`dirname $(pwd)`/tests

dir="$( cd "$( dirname "$0"  )" && pwd  )"
echo "xxxcvadsfadsafas"
echo $dir

function rungotest()
{
packages=$(go list ../... | grep -v -E 'vendor|tests'|grep $1) 
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
    #/usr/bin/docker run --rm -v $volumeDir:$gopath -e HOST_IP=$1 -e KUBE_MASTER_URL=$2 -e NODE_IP=$3 -e REGISTRY_BASE_URI=$4 -w $gopath $golangImage go test -v -cover -coverprofile=profile.tmp -coverpkg "$deps" $package
#    echo "/usr/bin/docker run --rm -v $volumeDir:$gopath --env-file env.cfg -w $gopath $golangImage go test -v -cover -coverprofile=profile.tmp -coverpkg "$deps" $package"
    echo "go test -v -cover -coverprofile=profile.tmp -coverpkg "$deps" $package"
    go test -v -cover -coverprofile=profile.tmp -coverpkg "$deps" $package
    if [ -f $volumeDir/profile.tmp ]
    then
        cat $volumeDir/profile.tmp | tail -n +2 >> $volumeDir/total.cov
        rm $volumeDir/profile.tmp
    fi

done
}
echo "mode: set" >$volumeDir/total.cov
rungotest apiserver 
rungotest tokenserver 
go tool cover -func=total.cov >> $volumeDir/total.temp
cov=`cat $dir/total.temp|grep "total"|grep -v -E 'NaN'|awk '{print $NF}'|cut -d "%" -f 1|tr -s [:space:]`
echo $cov > $dir/avaCov.cov
#return $cov

go tool cover -html=$dir/total.cov -o $dir/profile.html

