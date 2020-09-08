version: '2'
services:
  log:
    image: board_log:__version__
    restart: always
    volumes:
      - /var/log/board/:/var/log/docker/
      - /etc/localtime:/etc/localtime:ro
    networks:
      - board
    ports:
      - 1514:514
  db:
    image: board_db:__version__
    restart: always
    volumes:
      - /data/board/database:/var/lib/mysql
      - /etc/localtime:/etc/localtime:ro
    env_file:
      - ../config/db/env
    networks:
      - board
    depends_on:
      - log
    logging:
      driver: "syslog"
      options:  
        syslog-address: "tcp://127.0.0.1:1514"
        tag: "db"
  apiserver:
    image: board_apiserver:__version__
    restart: always
    volumes:
#     - ../../tools/swagger/vendors/swagger-ui-2.1.4/dist:/usr/bin/swagger:z
      - /data/board/cert:/cert:rw
      - ../config/apiserver/kubeconfig:/root/kubeconfig
      - /etc/board/cert:/etc/board/cert:rw
      - /etc/localtime:/etc/localtime:ro
    env_file:
      - ../config/apiserver/env
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
    image: board_tokenserver:__version__
    env_file:
      - ../config/tokenserver/env
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
    image: board_proxy:__version__
    networks:
      - board
    restart: always
    volumes:
      - ../config/proxy/nginx.conf:/etc/nginx/nginx.conf:z
#     - ../../src/ui/dist:/usr/share/nginx/html:z
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
    image: board_prometheus:__version__
    restart: always
    networks:
      - dvserver_net
    ports:
      - 9090:9090
    volumes:
      - ../config/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml
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
