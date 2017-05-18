#!/bin/bash
SCHEME=http
SERVER_IP=localhost

set -e

usage=$'Please set SERVER_IP in prepare-swagger.sh first. DO NOT use local
host or 127.0.0.1 for hostname, because this service needs to be accessed by external clients.'

while [ $# -gt 0 ]; do
        case $1 in
            --help)
            echo "$usage"
            exit 0;;
            *)
            echo "$usage"
            exit 1;;
        esac
        shift || true
done

# The SERVER_IP in prepare-swagger.sh has not been modified
if [ $SERVER_IP = localhost ]
then
        echo "$usage"
        exit 1
fi

cp ../../docs/swagger.yaml swagger.yaml

FILE="swagger.tar.gz"
if ! [ -f $FILE ]; then
    mkdir yaml vendors
    echo "Doing some clean up..."
    rm -f *.tar.gz
    echo "Downloading Swagger UI release package..."
    wget https://github.com/swagger-api/swagger-ui/archive/v2.1.4.tar.gz -O swagger.tar.gz
    echo "Untarring Swagger UI package to the static file path..."
    tar -C ./vendors -zxf swagger.tar.gz swagger-ui-2.1.4/dist
    echo "Executing some processes..."
    sed -i.bak 's/http:\/\/petstore\.swagger\.io\/v2\/swagger\.json/'$SCHEME':\/\/'$SERVER_IP'\/yaml\/swagger\.yaml/g' \
    ./vendors/swagger-ui-2.1.4/dist/index.html
    sed -i.bak '/jsonEditor: false,/a\        validatorUrl: null,' ./vendors/swagger-ui-2.1.4/dist/index.html
    mkdir -p ./yaml
    cp swagger.yaml ./yaml
    sed -i.bak 's/host: localhost/host: '$SERVER_IP'/g' ./yaml/swagger.yaml
    sed -i.bak 's/  \- http$/  \- '$SCHEME'/g' ./yaml/swagger.yaml
    echo "Finish preparation for the Swagger UI."

fi

docker run -d -p 80:80 -v "$PWD":/usr/local/apache2/htdocs/ httpd
echo "Start docker container for Swagger, please visit http://$SERVER_IP/vendors/swagger-ui-2.1.4/dist/index.html."
