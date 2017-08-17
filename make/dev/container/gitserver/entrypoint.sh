#!/bin/sh

# If there is some public key in keys folder
# then it copies its contain in authorized_keys file
if [ "$(ls -A /gitserver/keys/)" ]; then
  cd /home/git
  cat /gitserver/keys/*.pub > .ssh/authorized_keys
  chown -R git:git .ssh
  chmod 700 .ssh
  chmod -R 600 .ssh/*
fi

# Checking permissions and fixing SGID bit in repos folder
if [ "$(ls -A /gitserver/repos/)" ]; then
  cd /gitserver/repos
  chown -R git:git .
  chmod -R ug+rwX .
  find . -type d -exec chmod g+s '{}' +
fi

# -D flag avoids executing sshd as a daemon
/usr/sbin/sshd -D
