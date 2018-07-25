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


import libvirtapi

kmvNameList = ['kvm-1', 'kvm-2', 'kvm-3', 'kvm-4']
tmpkvmdir = '/tmp/kvm'

basedir = os.path.dirname(os.path.realpath(__file__))

credentialId = 'k-v-m-i-d'

kvmnodeUserName = 'root'
kvmnodePasswd = '123456a?'


def get_ip(ifname):
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    return socket.inet_ntoa(fcntl.ioctl(s.fileno(), 0x8915, struct.pack('256s', ifname[:15]))[20:24])


def createCredential(jenkinsMaster, nodeUsername, nodePasswd, credentialId):
    para = "--data-urlencode"
    post_url = '%s/credentials/store/system/domain/_/createCredentials ' %jenkinsMaster
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
    print (cmd)
    os.system(cmd)


def getNodeIp(nodename):
    macline = commands.getoutput('virsh dumpxml %s|grep mac|grep address' %nodename)
    mac = macline.split("'")[1]
    ipaddr = commands.getoutput('''arp -a|grep %s|cut -d '(' -f2|cut -d ')' -f1''' %mac)
    return ipaddr
def checkNode(jenkinsServer, usekvmnode):
    try:
        jenkinsServer.assert_node_exists(usekvmnode, exception_message='node %s does not exist.......')
    except jenkins.JenkinsException, e:
        print (e)
    nodeflag = jenkinsServer.node_exists(usekvmnode) 
    return nodeflag

def kvmIp(nodename):
    ipaddr = getNodeIp(nodename)
    while len(ipaddr)<7:
       # time.sleep(3)
        #ipaddr = getNodeIp()
        status,ipaddr=commands.getstatusoutput("%s/genIp.py %s" %(basedir,nodename))
    print ("==============:::" + ipaddr)
    return ipaddr
def addJenkinsNode(jenkinsMaster, nodename, hostport):
    ipaddr = get_ip('eno1')
    cid = credentialId
    server = jenkins.Jenkins(jenkinsMaster)
    flag = checkNode(server, nodename)
    
    if flag == False:
       
    
    #while flag == True:
    #    time.sleep(2)
    #    flag = checkNode(server, nodename)
    
        params = {
            'port': hostport,
            'username': 'juser',
            'credentialsId': cid,
            'host': ipaddr
        }
        print ("params: %s" %params)
        try:
            server.create_node(
                nodename,nodeDescription='add slave',
                remoteFS='/data',
                numExecutors=3,
                labels='slave',
                exclusive=False,
                launcher=jenkins.LAUNCHER_SSH,
                launcher_params=params)
        except Exception, e:
           print 'str(e):\t\t', str(e)
           print ("failed add jenkins node")
def checkKVM(flagFile):
    if os.path.isfile(flagFile):
        with open(flagFile, 'r') as f:
            lines = f.readlines()
            flag = lines[0].replace("\n", "")
            f.close
    print ("==%s==" %flag)
    return flag

def copyImage(kvmName):
    baseimage='/var/lib/libvirt/images/kvm.img'
    image = '/var/lib/libvirt/images/%s.img' %kvmName
    if os.path.exists(image):
        os.remove(image)
    else:
        print ('no image can be cleaned!!!')
    if os.path.exists(baseimage):
        shutil.copyfile(baseimage,image)
    else:
        print ('no base image can be copied!!!')

def diff(listA,listB):
    notused = list(set(listB).difference(set(listA)))
    return notused

def getdifflist():
    listB = kmvNameList
    listA = getlistofkvmDir()
    notusedkvmname = diff(listA,listB)
    return notusedkvmname

def getNumberofkvm():
    count = 0
    for fn in os.listdir(tmpkvmdir):
        count = count + 1
    return count
def getlistofkvmDir():
    for root, dirs, files in os.walk(tmpkvmdir):
        usingKvmList = files
    return usingKvmList
def getKvmName(kvmApiServer, projectName):
    url = '%s/register-job' %kvmApiServer
    data = {'job_name': projectName}
    res = requests.post(url, data=data)
    kvmName = res.text
    return kvmName
def triggerKvmJob(jenkinsMaster, jobName, kvmName):
    url = '%s/job/%s/buildWithParameters' % (jenkinsMaster, jobName)
    data = {'node_name': kvmName}
    res = requests.get(url, params=data)
    return res.text
def startKVM_1(jenkinsMaster, projectName, kvmApiServer):
    kvmName = getKvmName(kvmApiServer, projectName)
    print ("::::" + kvmName + ":::")
    while kvmName == 'FULL':
        time.sleep(3)
        print ('kvm is full, can not get a kvm name...')
        kvmName = getKvmName(kvmApiServer, projectName)

    usekvmname = kvmName
    print(usekvmname)

    conn = libvirtapi.createConnection()
    myDom = libvirtapi.getDomInfoByName(conn, usekvmname)
    libvirtapi.closeConnection(conn)
    
    if myDom == 1:
        print ('create kmv ...........................')
        copyImage(usekvmname)
        try:
            os.popen("virt-install --name %s --ram 2048 --disk path=/var/lib/libvirt/images/%s.img --import &\n\n\n" %(usekvmname, usekvmname))
        except:
            print('create kvm failed')
#    else:
#        print ('revert snapshot .........................')
#        try:
#            os.popen("virsh snapshot-revert %s %s" %(usekvmname, usekvmname))
#        except:
#            print('create kvm failed')
    print ('-----usekvmname---')
    print (usekvmname)
    hostport = '2000' + usekvmname.split('-')[-1]
    #confignat(usekvmname, hostport)
    addJenkinsNode(jenkinsMaster,usekvmname, hostport)

    if myDom == 1:
        try:
            os.popen('virsh snapshot-create-as %s %s' %(usekvmname, usekvmname))
        except:
            print('create kvm failed')
    triggerKvmJob(jenkinsMaster, projectName, usekvmname)

def confignat(usekvmname, hostport):
    kvmIpaddr = kvmIp(usekvmname)
    hostip = get_ip('eno1')
    kvmnodemask = '192.168.122.0/255.255.255.0'
    kvmnodegetway = '192.168.122.1'
    try:
        os.popen('iptables -A INPUT -p tcp --dport %s -j ACCEPT' %hostport)
        os.popen('iptables -t nat -A PREROUTING -d %s -p tcp -m tcp --dport %s -j DNAT --to-destination %s:22' %(hostip, hostport, kvmIpaddr))
        os.popen('iptables -t nat -A POSTROUTING -s %s -d %s -p tcp -m tcp --dport 22 -j SNAT --to-source %s' %(kvmnodemask, hostip, kvmnodegetway))
    except:
        print ('add iptables failed')
def startService():
    print (sta)
    os.system('systemctl start docker')
def main():
    jenkinsMaster = sys.argv[1]
    projectName = sys.argv[2]
    kvmApiServer = sys.argv[3]
    createCredential(jenkinsMaster, kvmnodeUserName, kvmnodePasswd, credentialId)
    sta = startKVM_1(jenkinsMaster, projectName, kvmApiServer)
    print ("add jenkins node ..................")
   

if __name__ == "__main__":
    main()
