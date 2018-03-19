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
gopath=/go/src/git/inspursoft/board/
golangImage=golang:1.8.3-alpine3.5
volumeDir=`dirname $(pwd)`

dir="$( cd "$( dirname "$0"  )" && pwd  )"

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
    echo "/usr/bin/docker run --rm -v $volumeDir:$gopath --env-file env.cfg -w $gopath $golangImage go test -v -cover -coverprofile=profile.tmp -coverpkg "$deps" $package"
    /usr/bin/docker run --rm -v $volumeDir:$gopath --env-file env.cfg -w $gopath $golangImage go test -v -cover -coverprofile=profile.tmp -coverpkg "$deps" $package

    if [ -f $volumeDir/profile.tmp ]
    then
        cat $volumeDir/profile.tmp | tail -n +2 >> profile.cov
        rm $volumeDir/profile.tmp
    fi

done
cp $dir/profile.cov $dir/$1".cov"
go tool cover -func=$1".cov" >> $1".temp"
cov=`cat $dir/$1".temp"|grep "total"|grep -v -E 'NaN'|awk '{print $NF}'|cut -d "%" -f 1|tr -s [:space:]`
echo $cov > $dir/$1".txt"
#return $cov
}
rungotest apiserver
cov1=`cat $dir/apiserver.txt`
rungotest tokenserver
cov2=`cat $dir/tokenserver.txt`

echo "--------------------"
echo $cov1
echo $cov2
echo "------------------"
add=$(echo $cov1+$cov2|bc)
averageCov=$(echo "scale=2;$add/2"|bc)
echo $averageCov>>$dir/avaCov.cov
go tool cover -html=profile.cov -o profile.html

