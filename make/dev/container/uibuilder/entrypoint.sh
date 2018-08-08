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
    rm -rf out-ngc/
    npm run aot
    cp src/rollup/main-aot.js out-ngc/src/main.js
    npm run rollup
    cp node_modules/@clr/ui/clr-ui.min.css dist/clr-ui.min.css
    cp node_modules/@clr/icons/clr-icons.min.css dist/clr-icons.min.css
    cp node_modules/inspurprism/prism-default.css dist/prism-default.css
    cp node_modules/es6-shim/es6-shim.min.js dist/es6-shim.min.js
    cp node_modules/echarts/dist/echarts.min.js dist/echarts.min.js
    cp node_modules/zone.js/dist/zone.min.js dist/zone.min.js
    cp node_modules/@webcomponents/custom-elements/custom-elements.min.js dist/custom-elements.min.js
    cp node_modules/@clr/icons/clr-icons.min.js dist/clr-icons.min.js
    cp src/styles.css dist/styles.css
    cp src/favicon.ico dist/favicon.ico
    cp src/rollup/index.html dist/index.html
    cp -R src/images dist/
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
