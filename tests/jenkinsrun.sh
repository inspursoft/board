BUILD_URL=$1
WORKSPACE=$2
head_repo_url=$3
head_branch=$4
base_repo_url=$5
base_branch=$6
comments_url=$7
JOB_URL=$8
JENKINS_URL=$9


totalLink=$BUILD_URL/TOTAL_REPORT/index.html
uiLink=$BUILD_URL/UI/index.html
consoleLink=$BUILD_URL/console
boardDir=$WORKSPACE/src/git/inspursoft
branchDir=`echo $head_repo_url|awk -F '/' '{print $NF}'|cut -d '.' -f 1`
workDir=$WORKSPACE
uiDir=$boardDir/$branchDir/src/ui


lastBuildCov=`curl $JOB_URL/lastSuccessfulBuild/COVERAGE/index.html|grep "%"|cut -d '>' -f 3|cut -d '%' -f 1`
lastUiBuildCov=`curl $JOB_URL/lastSuccessfulBuild/UICOVERAGE/index.html|grep "%"|cut -d '>' -f 3|cut -d '%' -f 1`

echo "--------------------------"
echo $lastBuildCov
echo "xxxxxxxxxxxxxxxxxxxxxxxxxx"


#make prepare
cd $boardDir/$branchDir
make prepare

#start up mysql docker container
#cp /home/backup/docker-compose.mysql.a.yml $boardDir/$branchDir/make/dev
cp $boardDir/$branchDir/tests/docker-compose.test.yml  $boardDir/$branchDir/make/dev
cp $boardDir/$branchDir/tests/ldap_test.ldif  $boardDir/$branchDir/make/dev
cd $boardDir/$branchDir/make/dev
docker-compose -f docker-compose.test.yml down -v
rm -rf /data/board
rm -rf /tmp/test-repos /tmp/test-keys
rm -f  /root/.ssh/known_hosts
docker-compose -f docker-compose.test.yml up -d


#docker-compose -f docker-compose.uibuilder.test.yml up 

export GOPATH=$workDir

cd $boardDir/$branchDir/tests


#cd $boardDir/board/tests

chmod +x *
envFile=$boardDir/$branchDir/tests/env.cfg
#make run
./run.sh $envFile

cp -r /home/tests/testresult.log /home/tests/coverage/ $uiDir
uiCoverage=`cat $uiDir/testresult.log |grep "Statements"|cut -d ":" -f 2|cut -d "%" -f 1|awk 'gsub(/^ *| *$/,"")'`

#cov=`cat $boardDir/$branchDir/tests/out.temp|grep "total"|awk '{print $NF}'|cut -d "%" -f 1|tr -s [:space:]`
covfile=$boardDir/$branchDir/tests/avaCov.cov


#echo "averageCov: " $averageCov



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

if [ ! -f $covfile ];then
pic="error.jpg"
uipic=`getFlag $lastUiBuildCov $uiCoverage`
cov="FAIL"
else
cov=`cat $covfile`"%"
add=`echo $cov+$uiCoverage|bc`
averageCov=`echo $add/2|bc`
echo "python genResult.py $WORKSPACE $cov $uiCoverage"
python genResult.py $WORKSPACE $cov $uiCoverage
pic=`getFlag $lastBuildCov $cov`
uipic=`getFlag $lastUiBuildCov $uiCoverage`
fi
echo $comments_url


echo "=================================================================="
info1="The test coverage for backend is "
commenturltmp=`echo $comments_url|sed 's/pulls/issues/g'|sed 's/inspursoft/api\/v1\/repos\/inspursoft/g'`
commentsurl=$commenturltmp"/comments"
serverresult=$cov
#commentsurl="http://10.110.18.40:10080/api/v1/repos/inspursoft/board/issues/1396/comments"
uiinfo=",The test coverage for frontend is "
consoleinfo=", check "
imageLink=$JENKINS_URL/userContent/$pic
uiImageLink=$JENKINS_URL/userContent/$uipic
imageuri=" <img src="$imageLink" width="20" height="20"> "
uiImageuri=" <img src="$uiImageLink" width="20" height="20"> "
uiImageuri=" <img src="$uiImageLink" width="20" height="20"> "
uiCov=" <a href=$uiLink>$uiCoverage </a>"
serverCovLink=" <a href=$totalLink>$serverresult</a>"
consoleuri=" <a href=$consoleLink> consolse log</a> "
imageLink=$JENKINS_URL/userContent/$pic
bodyinfo=$info1$serverCovLink$imageuri$uiinfo$uiImageuri$consoleinfo$consoleuri
bodyinfo=$info1$serverCovLink$imageuri$uiinfo$uiCov$uiImageuri$consoleinfo$consoleuri
b='-d { "body": "'$bodyinfo'"}'
cmd="curl -X POST \
  $commentsurl \
  -H 'Authorization: token 7c67f6a0d7967a152329c33ecaed3e1f93cb1e2d' \
  -H 'Content-Type: application/json' \
  '$b'"
echo $cmd
eval $cmd

echo "+++++++++++++++++=="
echo "comments_url	:"$comments_url
echo "info1		:"$info1
echo "uiinfo		:"$uiinfo
echo "imageuri		:"$imageuri
echo "uiImageuri	:"$uiImageuri
echo "uiCov		:"$uiCov
echo "serverCovLink	:"$serverCovLink
