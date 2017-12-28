#!/bin/bash
set -e

rm -rf /board_src/dist

cp -R /board_src/* /board_resource/.

cd /board_resource

cat package.json

npm install

ng build

cp -R dist /board_src
