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


totalLink=$BUILD_URL/TOTAL_REPORT/index.html
uiLink=$BUILD_URL/UI/index.html
consoleLink=$BUILD_URL/console
boardDir=$WORKSPACE/src/git/inspursoft
branchDir=`echo $head_repo_url|awk -F '/' '{print $NF}'|cut -d '.' -f 1`

uiDir=$boardDir/$branchDir/src/ui


lastBuildCov=`curl $JOB_URL/lastSuccessfulBuild/COVERAGE/index.html|grep "%"|cut -d '>' -f 3|cut -d '%' -f 1`
lastUiBuildCov=`curl $JOB_URL/lastSuccessfulBuild/UICOVERAGE/index.html|grep "%"|cut -d '>' -f 3|cut -d '%' -f 1`

echo "xxxxxxxxxxxxxxxxxxxxxxxxxx"
echo $lastBuildCov
echo "xxxxxxxxxxxxxxxxxxxxxxxxxx"


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


docker-compose -f docker-compose.uibuilder.test.yml up 

export GOPATH=$workDir

cd $boardDir/$branchDir/tests


#cd $boardDir/board/tests

chmod +x *
#make run
./run.sh $host_ip $kube_master_url $node_ip $registry_uri

#cp -r /home/tests/testresult.log /home/tests/coverage/ $uiDir
uiCoverage=`cat $uiDir/testresult.log |grep "Statements"|cut -d ":" -f 2|cut -d "%" -f 1|awk 'gsub(/^ *| *$/,"")'`

cov=`cat $boardDir/$branchDir/tests/out.temp|grep "total"|awk '{print $NF}'|cut -d "%" -f 1|tr -s [:space:]`

echo '==========================================='
echo $lastBuildCov
echo $cov
echo '==========================================='


add=`echo $cov+$uiCoverage|bc`
averageCov=`echo $add/2|bc`
echo "averageCov: " $averageCov

echo "python genResult.py $WORKSPACE $cov $uiCoverage"
python genResult.py $WORKSPACE $cov $uiCoverage


cp -r $uiDir/coverage $WORKSPACE/total


function getFlag()
{
lastC=$1
nowC=$2
ftmp=`echo "$lastC>$nowC"|bc `
if [ $ftmp -eq 1 ]; then
flag="down"
pic="error.jpg"
elif [ $ftmp -eq 0 ]; then
flag="eq"
pic="correct.jpg"
else
flag="up"
pic="correct.jpg"
fi
echo $pic
}

pic=`getFlag $lastBuildCov $cov`
uipic=`getFlag $lastUiBuildCov $uiCoverage`

echo $comments_url

tmp="?content=the%20API%20Server%20Coverage%20is%20"
serverCovLink="%20<a%20href=$totalLink>$cov%25</a>"
imageLink=$JENKINS_URL/userContent/$pic
uiImageLink=$JENKINS_URL/userContent/$uipic
image="%20<img%20src="$imageLink"%20width="20"%20height="20">%20"
uiImage="%20<img%20src="$uiImageLink"%20width="20"%20height="20">%20"
uiCov="%20UI%20Coverage%20is%20<a%20href=$uiLink>$uiCoverage%25</a>"
#f_comments_ur:=$comments_url$tmp$cov%25%20$totalLink
f_comments_url="$comments_url$tmp$serverCovLink$image$uiCov$uiImage%20check%20<a%20href=$consoleLink>console%20log</a>"

echo $f_comments_url


echo "curl --user Jenkins-10.110.18.40:123456 -X POST -H "content-type: application/x-www-form-urlencoded" $f_comments_url"
echo `curl --user 'Jenkins-10.110.18.40:123456' -X POST -H 'content-type: application/x-www-form-urlencoded' $f_comments_url`

