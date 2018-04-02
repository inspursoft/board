#! /bin/bash -e
cd /usr/share/jenkins/setuptools-38.5.1
python setup.py install
cd /usr/share/jenkins/python-jenkins-0.4.15
python setup.py install
cd /usr/share/jenkins/multi_key_dict-2.0.3
python setup.py install
