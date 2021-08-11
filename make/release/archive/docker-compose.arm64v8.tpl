version: '2'
services:
  log:
    image: openboard/board_log:__version__
    restart: always
    volumes:
      - /var/log/board/:/var/log/docker/
      - /etc/localtime:/etc/localtime:ro
    networks:
      - board
    ports:
      - 1514:514
  db:
    image: openboard/board_db:__version__
    restart: always
    volumes:
      - /data/board/database:/var/lib/mysql
      - /etc/localtime:/etc/localtime:ro
    env_file:
      - ../../config/db/env
    networks:
      - board
    depends_on:
      - log
    logging:
      driver: "syslog"
      options:  
        syslog-address: "tcp://127.0.0.1:1514"
        tag: "db"
  gogits:
    image: openboard/board_gogits:__version__
    restart: always
    volumes:
      - /data/board/gogits:/data:rw
      - ../../config/gogits/conf/app.ini:/tmp/conf/app.ini
      - /etc/localtime:/etc/localtime:ro
    ports:
      - "10022:22"
      - "10080:3000"
    networks:
      - board
    depends_on:
      - log
    logging:
      driver: "syslog"
      options:
        syslog-address: "tcp://127.0.0.1:1514"
        tag: "gogits"
  jenkins:
    image: openboard/board_jenkins:__version__
    restart: always
    networks:
      - board
    volumes:
      - /data/board/jenkins_home:/var/jenkins_home
      - ../../config/ssh_keys:/root/.ssh
      - /var/run/docker.sock:/var/run/docker.sock
      - /usr/bin/docker:/usr/bin/docker
      - /etc/localtime:/etc/localtime:ro
    env_file:
      - ../../config/jenkins/env
    ports:
      - 8888:8080
    depends_on:
      - log
    logging:
      driver: "syslog"
      options:
        syslog-address: "tcp://127.0.0.1:1514"
        tag: "jenkins"
  apiserver:
    image: openboard/board_apiserver:__version__
    restart: always
    volumes:
#     - ../../../tools/swagger/vendors/swagger-ui-2.1.4/dist:/usr/bin/swagger:z
      - /data/board/repos:/repos:rw
      - /data/board/keys:/keys:rw
      - /data/board/cert:/cert:rw
      - ../../config/apiserver/kubeconfig:/root/kubeconfig
      - /etc/board/cert:/etc/board/cert:rw
      - /etc/localtime:/etc/localtime:ro
    env_file:
      - ../../config/apiserver/env
    ports:
      - 8088:8088
    networks:
      - board
    links:
      - db
    depends_on:
      - log
    logging:
      driver: "syslog"
      options:
        syslog-address: "tcp://127.0.0.1:1514"
        tag: "apiserver"
  tokenserver:
    image: openboard/board_tokenserver:__version__
    env_file:
      - ../../config/tokenserver/env
    restart: always
    networks:
      - board
    depends_on:
      - log
    volumes:
      - /etc/localtime:/etc/localtime:ro
    logging:
      driver: "syslog"
      options:
        syslog-address: "tcp://127.0.0.1:1514"
        tag: "tokenserver"
  proxy:
    image: openboard/board_proxy:__version__
    networks:
      - board
    restart: always
    volumes:
      - ../../config/proxy/nginx.conf:/etc/nginx/nginx.conf:z
#     - ../../../src/ui/dist:/usr/share/nginx/html:z
      - /etc/localtime:/etc/localtime:ro
    ports: 
      - 80:80
      - 8080:8080
    links:
      - apiserver
    depends_on:
      - log
    logging:
      driver: "syslog"
      options:
        syslog-address: "tcp://127.0.0.1:1514"
        tag: "proxy"
  prometheus:
    image: openboard/board_prometheus:__version__
    restart: always
    networks:
      - dvserver_net
    ports:
      - 9090:9090
    volumes:
      - ../../config/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml
      - /etc/localtime:/etc/localtime:ro
    depends_on:
      - log
    logging:
      driver: "syslog"
      options:
        syslog-address: "tcp://127.0.0.1:1514"
        tag: "prometheus"
networks:
  board:
    external: true
