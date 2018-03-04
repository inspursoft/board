FROM jenkins:2.60.3

ENV JAVA_OPTS "-Djenkins.install.runSetupWizard=false \
               -Dhudson.model.User.allowNonExistentUser=true"

USER root

RUN sed -i 's/^root\:x\:0\:/root\:x\:0\:root\,jenkins/' /etc/group

COPY ./sources.list /tmp/sources.list
RUN cat /tmp/sources.list > /etc/apt/sources.list \
      && apt-get update \
      && apt-get install -y --allow-unauthenticated sudo libltdl7 curl\
      && rm -rf /var/lib/apt/lists/*
RUN echo "jenkins ALL=NOPASSWD: ALL" >> /etc/sudoers

COPY ./jobs /tmp/jobs
COPY ./python-jenkins-0.4.15 /usr/share/jenkins/python-jenkins-0.4.15
COPY ./setuptools-38.5.1 /usr/share/jenkins/setuptools-38.5.1
COPY ./multi_key_dict-2.0.3 /usr/share/jenkins/multi_key_dict-2.0.3
COPY ./pythonenv.sh /usr/share/jenkins/pythonenv.sh

RUN chmod +x /usr/share/jenkins/pythonenv.sh

COPY ./plugins.txt /usr/share/jenkins/plugins.txt

COPY ./addNode.py /usr/share/jenkins/addNode.py


RUN /usr/local/bin/plugins.sh /usr/share/jenkins/plugins.txt



COPY ./jenkins.sh /usr/local/bin/jenkins.sh
RUN chmod u+x /usr/local/bin/jenkins.sh

COPY ./addnode.sh /usr/share/jenkins/addnode.sh
RUN chmod u+x /usr/share/jenkins/addnode.sh
COPY ./init.sh /usr/share/jenkins/init.sh
RUN chmod u+x /usr/share/jenkins/init.sh

