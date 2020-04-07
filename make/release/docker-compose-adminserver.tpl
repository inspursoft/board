version: '2'
services:
  adminserver:
    image: board_adminserver:__version__
    restart: always
    volumes:
      - ./:/go/cfgfile
      - /data/board/secrets:/go/secrets
      - /var/run/docker.sock:/var/run/docker.sock
      - /data/board/database:/data/board/database
      - /data/board/ansible_k8s:/data/board/ansible_k8s
      - ./config:/data/board/Deploy/config
    env_file:
      - ../config/adminserver/env
    networks:
      - board
    ports:
      - 8081:8080
  proxy-adminserver:
    image: board_adminserver_proxy:__version__
    depends_on: 
      - adminserver
    restart: always
    ports:
      - 8082:80
    links:
      - adminserver
    volumes:
      - ./templates/proxy-adminserver/nginx.conf:/etc/nginx/nginx.conf:z
    networks:
      - board
networks:
  board:
    external: false
