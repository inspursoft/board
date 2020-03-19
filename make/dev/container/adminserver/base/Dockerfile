FROM golang:alpine

MAINTAINER zoltan@mullner.hu

RUN echo "===> Installing sudo to emulate normal OS behavior..."  && \
    apk --update add sudo                                         && \
    \
    \
    echo "===> Adding Python runtime..."  && \
    apk --update add python3 openssl ca-certificates    && \
    apk --update add --virtual build-dependencies \
                python3-dev libffi-dev openssl-dev build-base  && \
    pip3 install --upgrade pip cffi                            && \
    \
    \
    echo "===> Installing Ansible..."  && \
    pip3 install ansible==2.9.1         && \
    \
    \
    echo "===> Installing handy tools (not absolutely required)..."  && \
    pip3 install --upgrade pywinrm                  && \
    apk --update add sshpass openssh-client rsync  && \
    \
    \
    echo "===> Removing package list..."  && \
    apk del build-dependencies            && \
    rm -rf /var/cache/apk/*               && \
    \
    \
    echo "===> Adding hosts for convenience..."  && \
    mkdir -p /etc/ansible                        && \
    echo 'localhost' > /etc/ansible/hosts

RUN apk --no-cache --update add \
        bash \
        py-dnspython \
        py-boto \
        py-netaddr \
        bind-tools \
        html2text \
        php7 \
        php7-json \
        git \
        jq \
        curl

RUN pip3 install --no-cache-dir --upgrade yq

RUN pip3 install --no-cache-dir --upgrade mitogen

RUN mkdir -p /ansible/playbooks

WORKDIR /ansible/playbooks

VOLUME [ "/ansible/playbooks" ]

# default command: display Ansible version
CMD [ "ansible-playbook", "--version" ]
