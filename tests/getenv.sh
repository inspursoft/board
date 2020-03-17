#!/bin/sh


cd $boardDir

#git the para to tigger job

if [ "$action" == "push" ]; then
git clone --branch=$base_repo_branch $base_repo_clone_url
branchDir=`echo $base_repo_clone_url|awk -F '/' '{print $NF}'|cut -d '.' -f 1`

cd $boardDir/$branchDir
echo  "action is $action ....."

else
echo "action is $action...."
git clone --branch=$head_repo_branch $head_repo_clone_url
branchDir=`echo $head_repo_clone_url|awk -F '/' '{print $NF}'|cut -d '.' -f 1`

cd $boardDir/$branchDir

echo $base_repo_clone_url
git config --global user.email "you@example.com"
git remote add upstream $base_repo_clone_url
git fetch upstream
git checkout -b master-main --track upstream/$base_repo_branch
git merge $head_repo_branch 
fi

