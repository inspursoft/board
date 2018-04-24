FROM ubuntu:16.04

MAINTAINER wangkun_lc@inspur.com

COPY ./sources.list /etc/apt/sources.list

RUN  apt-get update \
  && apt-get install -y curl python python-pip git python-mysqldb mysql-client\
  && pip install -i https://pypi.tuna.tsinghua.edu.cn/simple alembic   \
  && mkdir -p /board-migration

WORKDIR /board-migration

COPY . .

RUN  chmod u+x ./run.sh \
  && mkdir ./backup

VOLUME [ "/board-migration/backup" ]

ENTRYPOINT ["./run.sh"]
