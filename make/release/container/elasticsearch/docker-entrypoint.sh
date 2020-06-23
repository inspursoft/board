#!/bin/bash

if [[ ! -f ${CERT_DIR}/bundle.zip ]]; then
    mkdir -p ${CERT_DIR}
    cat >${CERT_DIR}/instances.yaml <<EOF
instances:
  - name: elasticsearch
    dns:
      - elasticsearch
      - localhost
    ip:
      - 127.0.0.1
EOF
    /usr/share/elasticsearch/bin/elasticsearch-certutil cert --silent --pem --in ${CERT_DIR}/instances.yaml -out ${CERT_DIR}/bundle.zip;
    unzip ${CERT_DIR}/bundle.zip -d ${CERT_DIR}; 
fi;

exec /tini "--" "/usr/local/bin/docker-entrypoint.sh" "eswrapper"