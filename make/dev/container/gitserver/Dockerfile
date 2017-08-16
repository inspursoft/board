FROM alpine:3.4

MAINTAINER wangkun "wangkun_lc@inspur.com"

# "--no-cache" is new in Alpine 3.3 and it avoid using
# "--update + rm -rf /var/cache/apk/*" (to remove cache)

ENV http_proxy http://10.165.22.223:8118
ENV https_proxy http://10.165.22.223:8118

RUN apk add --no-cache \
  openssh \
  git

# Key generation on the server
RUN ssh-keygen -A

# SSH autorun
# RUN rc-update add sshd

WORKDIR /gitserver/

# -D flag avoids password generation
# -s flag changes user's shell
RUN mkdir /gitserver/keys \
  && adduser -D -s /usr/bin/git-shell git \
  && echo git:12345 | chpasswd \
  && mkdir /home/git/.ssh

# This is a login shell for SSH accounts to provide restricted Git access.
# It permits execution only of server-side Git commands implementing the
# pull/push functionality, plus custom commands present in a subdirectory
# named git-shell-commands in the user’s home directory.
# More info: https://git-scm.com/docs/git-shell
COPY make/dev/container/gitserver/git-shell-commands /home/git/git-shell-commands

# sshd_config file is edited for enable access key and disable access password
COPY make/dev/container/gitserver/sshd_config /etc/ssh/sshd_config
COPY make/dev/container/gitserver/entrypoint.sh entrypoint.sh
RUN chmod u+x entrypoint.sh


EXPOSE 22

VOLUME ["/gitserver/keys", "/gitserver/repos"]

CMD ["sh", "entrypoint.sh"]
