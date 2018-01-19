BUILD_URL=$1
WORKSPACE=$2
head_repo_url=$3
head_branch=$4
base_repo_url=$5
base_branch=$6
comments_url=$7


totalLink=$BUILD_URL/TOTAL_REPORT
consoleLink=$BUILD_URL/console

rm -rf $WORKSPACE/*
rm -rf /data/board/*
cd $WORKSPACE
boardDir=$WORKSPACE/src/git/inspursoft
mkdir -p $boardDir
mkdir -p $WORKSPACE/index
mkdir -p $WORKSPACE/tag
mkdir -p $WORKSPACE/total


#git the para to tigger job
echo "CONSOLE_LINKE=$consoleLink" > /var/log/errortrig.tmp
echo "COMMENT_URL=$comments_url" >> /var/log/errortrig.tmp


export GOPATH=$WORKSPACE

cd $boardDir

echo "-----------------------------------"
echo $head_repo_url
echo $head_branch
echo $base_repo_url
echo $base_branch
echo $comments_url
echo "+-----------------------------------"

git clone --branch=$head_branch $head_repo_url

branchDir=`echo $head_repo_url|awk -F '/' '{print $NF}'|cut -d '.' -f 1`

cd $boardDir/$branchDir

echo $base_repo_url
git remote add upstream $base_repo_url
git fetch upstream
git checkout -b master-main --track upstream/$base_branch
git merge $head_branch 

export TAG=`git describe --tags`
echo VVV=$TAG > /var/log/tag.tmp



