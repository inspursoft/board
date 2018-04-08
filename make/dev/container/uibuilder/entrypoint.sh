#!/bin/bash
set -e
cp -R /board_src/* /board_resource/.
cd /board_resource
cat package.json
echo -e "MODEL is:${MODEL}\n"
if [ ${MODEL} = "prod" ];then
 echo -e "Begin execute prod"
 rm -rf /board_src/dist
 npm run aot
 npm run mainjs
 npm run rollup
 npm run copyfiles
 cp -R dist /board_src
 echo -e "End execute prod"
elif [ ${MODEL} = "test" ];then
 echo -e "Begin execute test"
 rm -rf /board_src/coverage
 npm test > testresult.log
 cat testresult.log
 cp testresult.log /board_src
 cp -R coverage /board_src
 echo -e "End execute test"
else
 echo -e "Begin execute build"
 rm -rf /board_src/dist
 ng build
 cp -R dist /board_src
 echo -e "End execute build"
fi
