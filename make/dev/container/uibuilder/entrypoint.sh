#!/bin/bash
set -e

cd /board_src
rm -rf dist/*

cat ./package.json

npm install

ng build
