BUILD_URL=$1
WORKSPACE=$2
head_repo_url=$3
head_branch=$4
base_repo_url=$5
base_branch=$6
comments_url=$7
host_ip=$8
kube_master_url=$9
shift 9
node_ip=$1
registry_uri=$2
JOB_URL=$3
JENKINS_URL=$4


totalLink=$BUILD_URL/TOTAL_REPORT
consoleLink=$BUILD_URL/console
boardDir=$WORKSPACE/src/git/inspursoft
branchDir=`echo $head_repo_url|awk -F '/' '{print $NF}'|cut -d '.' -f 1`


echo "--------------------------------"
echo "curl $JOB_URL/COVERAGE/index.html|grep "%"|cut -d '>' -f 3|cut -d '%' -f 1"
echo "--------------------------------"

lastBuildCov=`curl $JOB_URL/lastSuccessfulBuild/COVERAGE/index.html|grep "%"|cut -d '>' -f 3|cut -d '%' -f 1`

echo "lastBuildCov-----------------------------"
echo $lastBuildCov
echo "--------------------------------"


#make prepare
cd $boardDir/$branchDir
make prepare

rm -rf $boardDir/$branchDir/src/collector/cmd/app/appflag_test.go

#start up mysql docker container
#cp /home/backup/docker-compose.mysql.a.yml $boardDir/$branchDir/make/dev
cp $boardDir/$branchDir/tests/docker-compose.test.yml  $boardDir/$branchDir/make/dev
cp $boardDir/$branchDir/tests/ldap_test.ldif  $boardDir/$branchDir/make/dev
cd $boardDir/$branchDir/make/dev
docker-compose -f docker-compose.test.yml down -v
docker-compose -f docker-compose.test.yml up -d


export GOPATH=$workDir
cd $boardDir/$branchDir/tests


cd $boardDir/board/tests

chmod +x *
#make run
./run.sh $host_ip $kube_master_url $node_ip $registry_uri

cov=`python genResult.py /home|grep "the cover"|cut -d ":" -f 2|cut -d "%" -f 1`

echo '==========================================='
echo $lastBuildCov
echo $cov
echo '==========================================='

ftmp=`echo "$lastBuildCov>$cov"|bc `
echo "------------------------"
echo $ftmp
echo "------------------------"
if [ $ftmp -eq 1 ]; then
flag="down"
pic="error.jpg"
elif [ $ftmp -eq 0 ]; then
flag="eq"
pic="correct.jpg"
else
flang="up"
pic="correct.jpg"
fi
echo $flag

python genResult.py $WORKSPACE

echo $comments_url

tmp="?content=the%20coverage%20is%20"
imageLink=$JENKINS_URL/userContent/$pic
image="%20<img%20src="$imageLink"%20width="20"%20height="20">%20"
f_comments_url="$comments_url$tmp%20<a%20href=$totalLink>$cov%25</a>$image%20check%20<a%20href=$consoleLink>console%20log</a>"

echo $f_comments_url


echo "curl --user Jenkins-10.110.18.40:123456 -X POST -H "content-type: application/x-www-form-urlencoded" $f_comments_url"
echo `curl --user 'Jenkins-10.110.18.40:123456' -X POST -H 'content-type: application/x-www-form-urlencoded' $f_comments_url`

