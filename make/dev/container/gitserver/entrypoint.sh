#!/bin/sh

# If there is some public key in keys folder
# then it copies its contain in authorized_keys file
if [ "$(ls -A /gitserver/keys/)" ]; then
  cd /root
  cat /gitserver/keys/*.pub > .ssh/authorized_keys
fi

# -D flag avoids executing sshd as a daemon
/usr/sbin/sshd -D
