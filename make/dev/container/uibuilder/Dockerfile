FROM node:7.5.0

RUN mkdir -p /board_resource
RUN mkdir -p /board_src

COPY src/ui/package.json /board_resource
COPY make/dev/container/uibuilder/entrypoint.sh /

WORKDIR /board_resource

RUN npm install -g @angular/cli  && \
    npm install && \
    chmod u+x /entrypoint.sh
VOLUME ["/board_src"]
