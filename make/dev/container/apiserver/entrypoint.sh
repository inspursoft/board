#!/bin/bash

rm -rf /root/.ssh/*

ssh-keygen -t rsa -f /root/.ssh/id_rsa -q -N ""
ssh-keyscan gitserver >> /root/.ssh/known_hosts

./apiserver
