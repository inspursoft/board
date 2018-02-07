#!/bin/bash
set -e

rm -rf /board_src/coverage

cp -R /board_src/* /board_resource/.

cd /board_resource

cat package.json

npm install

npm test > testresult.log

cp testresult.log /board_src

cp -R coverage /board_src
