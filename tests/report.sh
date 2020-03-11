#!/bin/sh

if [ $? -ne 0 ]; then
    result='failed'
else
    result='pass'
fi


consoleLink=$jenkins_master_url/job/$group_name/$build_id/console
last_build_cov=`echo $last_coverage|cut -d ":" -f 2`
last_ui_cov=`echo $last_coverage|cut -d ":" -f 1`
echo "--------------------------"
echo $build_id
echo $action

#make prepare

chmod +x *

covfile=$boardDir/$base_repo_name/tests/avaCov.cov
build_cov=`cat $boardDir/$base_repo_name/tests/avaCov.cov`
ui_cov=`cat $boardDir/$base_repo_name/src/ui/testresult.log |grep "Statements"|cut -d ":" -f 2|cut -d "%" -f 1|awk 'gsub(/^ *| *$/,"")'`

if [ "$action" == "pull_request" ]; then
coverage_build_html=$boardDir/$base_repo_name/tests/profile.html
coverage_ui_html=$boardDir/$base_repo_name/src/ui/coverage/index.html
coverage_ui_tar=$boardDir/$base_repo_name/src/ui/coverage.tar
echo $coverage_ui_tar
cd $boardDir/$base_repo_name/src/ui/
echo $boardDir/$base_repo_name/src/ui/
tar cvf coverage.tar coverage

echo "push to register======================="
echo "gogs_url:		$gogs_url"
echo "group_name:	$group_name"
echo "full_name:	$full_name"
echo "username:		$username"
echo "cov_num:		$cov_num"

for postfile in $coverage_build_html $coverage_ui_tar;
do
command="curl -X POST \
  '$gogs_url/upload?full_name=$full_name&build_number=$build_id' \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: multipart/form-data' \
  -F 'upload=@$postfile'
"
echo $command
eval $command
done
coverage_build_html_path="$gogs_url/results/$full_name/$build_id/profile.html"
coverage_ui_html_path="$gogs_url/results/$full_name/$build_id/coverage/index.html"

echo $coverage_build_html_path
echo $coverage_ui_html_path

commit_cov_num="curl -X PUT \
  '$gogs_url/config?group_name=$group_name&username=$username' \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -d '{
                \"config_key\": \"last_coverage\",
                \"config_val\": \"$build_cov:$ui_cov\"
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


echo "last_build_cov:$last_build_cov;build_cov:$build_cov;"
pic=`getFlag $last_build_cov $build_cov`
echo $comment_url

uiPic=`getFlag $last_ui_cov $ui_cov`

fi

echo "=================================================================="
info1="The test coverage for backend is "
serverresult=$build_cov
consoleinfo=", check "
serverCovLink=" <a href=$coverage_build_html_path>$serverresult%</a>"
uiCovLink=" <a href=$coverage_ui_html_path>$ui_cov%</a>"
console_url=" <a href=$consoleLink> console log</a> "
build_image_link=$gogs_url/results/pic/$pic
ui_image_link=$gogs_url/results/pic/$uiPic
build_image_url=" <img src="$build_image_link" width="20" height="20"> "
ui_image_url=" <img src="$build_image_link" width="20" height="20"> "
bodyinfo=$info1$serverCovLink$imageu$uiinfo$uiImageuri$consoleinfo$consoleuri
bodyinfo=$info1$serverCovLink$imageuri$uiinfo$uiCov$uiImageuri$consoleinfo$consoleuri
bodyinfo=$info1$serverCovLink$imageuri$uiinfo$uiCov$uiImageuri$consoleinfo$consoleuri
bodyinfo=$info1$serverCovLink$build_image_url$uiCovLink$ui_image_url$consoleinfo$console_url
#bodyinfo=$info1
b='-d { "body": "'$bodyinfo'"}'
echo $b
cmd="curl -X POST \
  $comment_url \
  -H 'Authorization: token $access_token_jenkins' \
  -H 'Content-Type: application/json' \
  '$b'"

echo "---------command--------"
echo $cmd
eval $cmd

echo "+++++++++++++++++++"
echo "comment_url	:"$comment_url
echo "info1		:"$info1
echo "uiinfo		:"$uiinfo

elif [ "$action" == "push" ]; then
cmd="curl -X POST '$gogs_url/commit-report' -H 'Content-Type: application/json' -d '{ \"commit_id\": \"${commit_id}\", \"report\": \"$result|${jenkins_master_url}/job/${group_name}/${build_number}/console\"}'"

echo $cmd
eval $cmd 
fi
