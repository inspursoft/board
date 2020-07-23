FROM python:3.6-slim-jessie
MAINTAINER wangkun_lc@inspur.com
RUN mkdir /app

COPY app /app

WORKDIR /app
RUN pip install -i http://mirrors.aliyun.com/pypi/simple/ --trusted-host mirrors.aliyun.com -e .

VOLUME [ "/app/instance", "/app/env" ]

CMD [ "python", "action/perform.py"]
