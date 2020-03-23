version: '2'
services:
  adminserver:
    image: board_adminserver:__version__
    restart: always
    volumes:
      - ../:/go/cfgfile
      - /data/board/database:/data/board/database
      - /data/board/ansible_k8s:/data/board/ansible_k8s
    networks:
      - adms
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
      - ../templates/proxy-adminserver/nginx.conf:/etc/nginx/nginx.conf:z
    networks:
      - adms
networks:
  adms:
    external: false
