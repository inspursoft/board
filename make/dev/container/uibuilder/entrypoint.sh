#!/bin/bash
set -e

cd /board_src
rm -rf dist/*

#Check if node_modules directory existing
if [ ! -d "./node_modules" ]; then
  mv /board_resource/node_modules ./
fi

cat ./package.json

npm install

ng build
