BUILD_URL=$1
WORKSPACE=$2
head_repo_url=$3
head_branch=$4
base_repo_url=$5
base_branch=$6
comments_url=$7


totalLink=$BUILD_URL/TOTAL_REPORT
consoleLink=$BUILD_URL/console
boardDir=$WORKSPACE/src/git/inspursoft
branchDir=`echo $head_repo_url|awk -F '/' '{print $NF}'|cut -d '.' -f 1`

#make prepare
cd $boardDir/$branchDir
make prepare

rm -rf $boardDir/$branchDir/src/collector/cmd/app/appflag_test.go

#start up mysql docker container
cp /home/backup/docker-compose.mysql.a.yml $boardDir/$branchDir/make/dev
cd  $boardDir/$branchDir/make/dev
docker-compose -f docker-compose.mysql.a.yml down -v
docker-compose -f docker-compose.mysql.a.yml up -d


yes|cp /home/jenkinsRun/user_test.go $boardDir/$branchDir/src/apiserver/service/user_test.go
yes|cp /home/jenkinsRun/collectordao_test.go $boardDir/$branchDir/src/collector/dao
yes|cp /home/jenkinsRun/init_test.go $boardDir/$branchDir/src/collector/service/collect
yes|cp /home/jenkinsRun/dashboard_test.go $boardDir/$branchDir/src/common/dao

yes|cp /home/backup/genResult.py $boardDir/$branchDir/tests
yes|cp /home/jenkinsRun/run.sh $boardDir/$branchDir/tests


export GOPATH=$workDir
cd $boardDir/$branchDir/tests


cd $boardDir/board/tests

chmod +x *
make run

cov=`python genResult.py /home|grep "the cover"|cut -d ":" -f 2|cut -d "%" -f 1`

python genResult.py $WORKSPACE

echo $comments_url

tmp="?content=the%20coverage%20is%20"
#f_comments_url=$comments_url$tmp$cov%25%20$totalLink
f_comments_url="$comments_url$tmp%20<a%20href=$totalLink>$cov%25</a>%20check%20<a%20href=$consoleLink>console%20log</a>"

echo $f_comments_url


echo "curl --user Jenkins-10.110.18.40:123456 -X POST -H "content-type: application/x-www-form-urlencoded" $f_comments_url"


