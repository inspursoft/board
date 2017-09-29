#!/bin/bash

#mkdir /root/.ssh
ssh-keyscan gitserver > /root/.ssh/known_hosts

./apiserver
