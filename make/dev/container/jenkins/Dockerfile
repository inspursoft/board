FROM jenkins/jenkins:2.60.3

ENV JAVA_OPTS "-Djenkins.install.runSetupWizard=false \
               -Dhudson.model.User.allowNonExistentUser=true"

USER root

RUN sed -i 's/^root\:x\:0\:/root\:x\:0\:root\,jenkins/' /etc/group

COPY sources.list /tmp/sources.list
RUN cat /tmp/sources.list > /etc/apt/sources.list \
      && apt-get update \
      && apt-get install -y --allow-unauthenticated sudo libltdl7\
      && rm -rf /var/lib/apt/lists/*
RUN echo "jenkins ALL=NOPASSWD: ALL" >> /etc/sudoers

COPY jobs /tmp/jobs

COPY plugins.txt /usr/share/jenkins/plugins.txt

RUN /usr/local/bin/plugins.sh /usr/share/jenkins/plugins.txt

COPY jenkins.sh /usr/local/bin/jenkins.sh
RUN chmod u+x /usr/local/bin/jenkins.sh
