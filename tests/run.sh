# listDeps lists packages referenced by package in $1, 
# excluding golang standard library and packages in 
# direcotry vendor

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
echo "mode: set" >profile.cov

# set envirnment
deps=""
gop=/go/src/git/inspursoft/board/
goversion=golang:1.8.3-alpine3.5
dock=/usr/bin/docker
var=`pwd`
PWD=`dirname $var`

packages=$(go list ../... | grep -v -E 'vendor|tests')
for package in $packages
do
    listDeps $package

    echo "DEBUG: testing package $package"
    echo "$deps"
    
    echo "---------------------------------------"
    echo $deps
    echo $package
    echo "+++++++++++++++++++++++++++++++++++++++"
    
    #go env used docker container
    echo "$dock run --rm -v $PWD:$gop -w $gop $goversion go test -v -cover -coverprofile=profile.tmp -coverpkg "$deps" $package"
    /usr/bin/docker run --rm -v $PWD:$gop -w $gop $goversion go test -v -cover -coverprofile=profile.tmp -coverpkg "$deps" $package

    if [ -f $PWD/profile.tmp ]
    then
        cat $PWD/profile.tmp | tail -n +2 >> profile.cov
        rm $PWD/profile.tmp
     fi
done
go tool cover -func=profile.cov > out.temp
go tool cover -html=profile.cov -o profile.html

