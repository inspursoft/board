source /root/env.cfg
build_number=$1
WORKSPACE=$2

consoleLink=$jenkins_master_url/job/$group_name/$build_number/console
boardDir=$WORKSPACE/src/git/inspursoft
branchDir=`echo $base_repo_clone_url|awk -F '/' '{print $NF}'|cut -d '.' -f 1`
workDir=$WORKSPACE

last_build_cov=$last_coverage
echo "--------------------------"
echo $lastBuildCov
echo $build_number
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

if [ "$action" == "pull_request" ]; then
#cov=`cat $boardDir/$branchDir/tests/out.temp|grep "total"|awk '{print $NF}'|cut -d "%" -f 1|tr -s [:space:]`
covfile=$boardDir/$branchDir/tests/avaCov.cov
coverage_file_html=$boardDir/$branchDir/tests/profile.html
build_cov=`cat $boardDir/$branchDir/tests/avaCov.cov`


echo "push to register======================="
echo "gogs_url:		$gogs_url"
echo "group_name:	$group_name"
echo "full_name:	$full_name"
echo "username:		$username"
echo "cov_num:		$cov_num"
command="curl -X POST \
  '$gogs_url/upload?full_name=$full_name&build_number=$build_number' \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: multipart/form-data' \
  -F 'upload=@$coverage_file_html'
"
echo $command
eval $command
coverage_file_html_path="$gogs_url/results/$full_name/$build_number/profile.html"

echo $coverage_file_html_path

commit_cov_num="curl -X PUT \
  '$gogs_url/config?group_name=$group_name&username=$username' \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -d '{
                \"config_key\": \"last_coverage\",
                \"config_val\": \"$build_cov\"
}'"

echo $commit_cov_num
eval $commit_cov_num

echo "======================="

#echo "averageCov: " $averageCov

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
cov="FAIL"
else

pic=`getFlag $last_build_cov $build_cov`
fi
echo $comment_url


echo "=================================================================="
info1="The test coverage for backend is "
serverresult=$build_cov
consoleinfo=", check "
uiCov=" <a href=$uiLink>$uiCoverage </a>"
serverCovLink=" <a href=$coverage_file_html_path>$serverresult</a>"
console_url=" <a href=$consoleLink> consolse log</a> "
imageLink=$gogs_url/results/pic/$pic
image_url=" <img src="$imageLink" width="20" height="20"> "
bodyinfo=$info1$serverCovLink$imageu$uiinfo$uiImageuri$consoleinfo$consoleuri
bodyinfo=$info1$serverCovLink$imageuri$uiinfo$uiCov$uiImageuri$consoleinfo$consoleuri
bodyinfo=$info1$serverCovLink$imageuri$uiinfo$uiCov$uiImageuri$consoleinfo$consoleuri
bodyinfo=$info1$serverCovLink$image_url$consoleinfo$console_url
#bodyinfo=$info1
b='-d { "body": "'$bodyinfo'"}'
echo $b
cmd="curl -X POST \
  $comment_url \
  -H 'Authorization: token $access_token' \
  -H 'Content-Type: application/json' \
  '$b'"
echo $cmd
eval $cmd

echo "+++++++++++++++++++"
echo "comment_url	:"$comment_url
echo "info1		:"$info1
echo "uiinfo		:"$uiinfo
fi
