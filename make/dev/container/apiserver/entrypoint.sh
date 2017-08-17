#!/bin/bash

ssh-keyscan gitserver >> /root/.ssh/known_hosts

./apiserver
