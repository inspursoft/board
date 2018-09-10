import os
import sys
import commands
import jenkins
import json
import time
import shutil
import requests
import socket
import fcntl
import struct
import re

import libvirtapi


basedir = os.path.dirname(os.path.realpath(__file__))
credentialId = 'k-v-m-i-d'
kvmnodeusername = 'root'
kvmnodepassword = '123456a?'

def createCredential(jenkinsmasterurl, nodeUsername, nodePasswd, credentialId):
    para = "--data-urlencode"
    post_url = '%s/credentials/store/system/domain/_/createCredentials ' %jenkinsmasterurl
    credent = {
    "scope": "GLOBAL",
    "id": credentialId,
    "username": nodeUsername,
    "password": nodePasswd,
    "description": "auto added",
    "$class": "com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl"}

    data = json.dumps(credent)
    cmd = '''curl -X POST %s %s 'json={"":"0","credentials":%s}'
          ''' %(post_url, para, data)
    print cmd
    os.system(cmd)


def getNodeIp(nodename):
    macline = commands.getoutput('virsh dumpxml %s|grep mac|grep address' %nodename)
    mac = macline.split("'")[1]
    ipaddr = commands.getoutput('''arp -a|grep %s|cut -d '(' -f2|cut -d ')' -f1''' %mac)
    return ipaddr

def getNodeStatus(jenkinsServer, usekvmnode):
    try:
        nodes = jenkinsServer.get_nodes()
    except jenkins.JenkinsException, e:
        print "Failed to get node, err: %s" % e

    print "Current allocated KVM name is: %s" % (usekvmnode,)
    for tmp in nodes:
        print (tmp)
        name = tmp['name']
        status = tmp['offline']

        if name == usekvmnode:
            break
    return status

def checkNode(jenkinsServer, usekvmnode):
    try:
        jenkinsServer.assert_node_exists(usekvmnode, exception_message='node %s does not exist.......')
    except jenkins.JenkinsException, err:
        print err
    nodeflag = jenkinsServer.node_exists(usekvmnode) 
    return nodeflag

def triggerJob(jenkinsmasterurl, jobname, kvmname, branch):
    url = '%s/job/%s/buildWithParameters' % (jenkinsmasterurl, jobname)
    data = {'node_name': kvmname, 'branch': branch}
    res = requests.get(url, params=data)
    return res.text

def addJenkinsNode(jenkinsmasterurl, nodename, jenkinsnodeip, hostport):
    cid = credentialId
    server = jenkins.Jenkins(jenkinsmasterurl)
    flag = checkNode(server, nodename)
    print "Node existing status: %s" % flag
    if flag == False:
        params = {
            'port': hostport,
            'username': 'juser',
            'credentialsId': cid,
            'host': jenkinsnodeip
        }
        print "params: %s" % params
        try:
            print "Adding %s as Jenkins node ..." % nodename
            server.create_node(
                nodename,nodeDescription='Added node: %s' % nodename,
                remoteFS='/data',
                numExecutors=3,
                labels=nodename,
                exclusive=False,
                launcher=jenkins.LAUNCHER_SSH,
                launcher_params=params)
        except Exception, err:
           print "Failed to add jenkins node, err: %s" % err
    nodestatus = getNodeStatus(server, nodename)
    print "Current node status is: %s" % nodestatus
    while nodestatus == True:
        time.sleep(3)
        nodestatus = getNodeStatus(server, nodename)

def getKVMName(kvmregistryurl, projectname):
    url = '%s/register-job' %kvmregistryurl
    data = {'job_name': projectname}
    res = requests.post(url, data=data)
    kvmname = res.text
    return kvmname

def initKVM(jenkinsmasterurl, projectname, jenkinsnodeip, kvmregistryurl, branch):
    kvmName = getKVMName(kvmregistryurl, projectname)
    print "::::" + kvmName + ":::"
    while kvmName == 'FULL':
        time.sleep(3)
        print "No available KVM currently, please wait for release..."
        kvmName = getKVMName(kvmregistryurl, projectname)
    usekvmname = kvmName
    try:
        os.popen("virsh snapshot-revert %s %s" %(usekvmname, usekvmname))
    except Exception, e:
        print "Failed to create KVM failed, err: %s" % e
    print "Initialized KVM: %s as Jenkins node" % usekvmname
    hostport = '2000' + usekvmname[-1]
    addJenkinsNode(jenkinsmasterurl, usekvmname, jenkinsnodeip, hostport)
    return usekvmname

def main():
    jenkinsmasterurl = sys.argv[1]
    jobname = sys.argv[2]
    branch = sys.argv[3]
    jenkinsnodeip = sys.argv[4]
    kvmregistryport = sys.argv[5]
    kvmregistryurl = "http://%s:%s" % (jenkinsnodeip, kvmregistryport)
    createCredential(jenkinsmasterurl, kvmnodeusername, kvmnodepassword, credentialId)
    kvmname = initKVM(jenkinsmasterurl, jobname, jenkinsnodeip, kvmregistryurl, branch)
    triggerJob(jenkinsmasterurl, jobname, kvmname, branch)
   
if __name__ == "__main__":
    main()