FROM node:11.13.0

RUN mkdir -p /board_resource

COPY src/ui/package.json /board_resource
COPY make/dev/container/uibuilder/entrypoint.sh /

WORKDIR /board_resource

RUN mkdir -p /board_src \ 
    && npm install -g @angular/cli@latest \
    && npm install \
    && chmod u+x /entrypoint.sh
VOLUME ["/board_src"]
