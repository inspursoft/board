version: '2'
services:
  adminserver:
    image: board_adminserver:__version__
    restart: always
    volumes:
      - ../:/go/cfgfile
      - /data/board/secrets:/go/secrets
      - /var/run/docker.sock:/var/run/docker.sock
      - /data/adminserver/database:/data/adminserver/database
      - /data/board/ansible_k8s:/data/board/ansible_k8s
      - ../config:/data/board/make/config
    env_file:
      - ./env
    networks:
      - board
    ports:
      - 8081:8080
  proxy-adminserver:
    image: board_proxy_adminserver:__version__
    depends_on: 
      - adminserver
    restart: always
    ports:
      - 8082:80
    links:
      - adminserver
    volumes:
      - ../templates/proxy-adminserver/nginx.conf:/etc/nginx/nginx.conf:z
    networks:
      - board
networks:
  board:
    external: false
