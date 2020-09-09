#! /bin/bash -e
cd /usr/share/jenkins
tar zxvf python-jenkins-0.4.15.tar.gz
tar zxvf multi_key_dict-2.0.3.tar.gz
unzip -o setuptools-38.5.1.zip
cd /usr/share/jenkins/setuptools-38.5.1
python setup.py install
cd /usr/share/jenkins/python-jenkins-0.4.15
python setup.py install
cd /usr/share/jenkins/multi_key_dict-2.0.3
python setup.py install
