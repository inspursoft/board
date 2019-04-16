#!/bin/bash
set -e
cp -R /board_src/* /board_resource/.
cd /board_resource
cat package.json
echo -e "Current mode is:${MODE}\n"

function prod(){
    echo -e "Begin executing prod"
    rm -rf /board_src/dist
    rm -rf dist/
    npm run prod
    cp -R dist /board_src
    echo -e "End executing prod"
}

function test(){
    echo -e "Begin executing test"
    rm -rf converage/
    rm -rf /board_src/coverage
    npm test > testresult.log
    cat testresult.log
    cp testresult.log /board_src
    cp -R coverage /board_src
    echo -e "End executing test"
}

function build(){
    echo -e "Begin executing build"
    rm -rf dist/
    rm -rf /board_src/dist
    ng build
    cp -R dist /board_src
    echo -e "End executing build"
}

case ${MODE} in
    prod)
        prod
        exit 0
    ;;
    test)
        test
        exit 0
    ;;
    *)
        build
        exit 0
    ;;
esac
