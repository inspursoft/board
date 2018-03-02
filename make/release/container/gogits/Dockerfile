FROM 10.110.18.177/board/gogs-custom:v1.2

COPY make/config/gogits/conf/app.ini /tmp/app.ini
COPY make/release/container/gogits/setup /app/gogs/docker/s6/gogs/setup
CMD ["/bin/s6-svscan", "/app/gogs/docker/s6/"]