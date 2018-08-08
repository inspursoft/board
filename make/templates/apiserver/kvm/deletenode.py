import os
import sys
import commands
import jenkins
import time
import requests

def deleteJenkinsNode(jenkinsMaster, projectName, kvmApiServer):
    #cid = '3ea2a0cd-b611-433b-9998-d2d8882f97b2'
    nodename = getKvmName(kvmApiServer, projectName)
    server = jenkins.Jenkins(jenkinsMaster)
    try:
        server.delete_node(nodename)
    except:
        print ('can not delete node')
def releaseKvmNode(kvmApiServer, projectName):
    url = '%s/release-node' %kvmApiServer
    data = {'job_name': projectName}
    res = requests.get(url, params=data)
    print (res.text)
def getKvmName(kvmApiServer, projectName):
    url = '%s/register-job' %kvmApiServer
    data = {'job_name': projectName}
    res = requests.post(url, data=data)
    kvmName = res.text
    return kvmName


def main():
    jenkinsMaster = sys.argv[1]
    projectName = sys.argv[2]
    kvmApiServer = sys.argv[3]
    deleteJenkinsNode(jenkinsMaster, projectName, kvmApiServer)
    releaseKvmNode(kvmApiServer, projectName)
   

if __name__ == "__main__":
    main()
