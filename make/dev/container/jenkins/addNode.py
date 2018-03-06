'''addNode.py $jenkins_home'''
import os,sys,re
import jenkins
import json
import urllib
import urllib2
import time


def getCredentialId(jenkins_home):
    credentialFile = os.path.join(jenkins_home,'credentials.xml')
    with open(credentialFile) as f:
        lines = f.readlines()
        for line in lines:
            m = (re.search('''<id>(.*)</id>''', line))
            if m is not None:
                cid = (m.group(1))
                print (cid)
                return cid

def curl(jenkinsMaster):

    post_url = '%s/credentials/store/system/domain/_/createCredentials' %jenkinsMaster

    postData = {
        "": "0",
        "credentials": {
        "scope": "GLOBAL",
        "username": "root",
        "privateKeySource": {"value": "2", "stapler-class": "com.cloudbees.jenkins.plugins.sshcredentials.impl.BasicSSHUserPrivateKey$UsersPrivateKeySource"},
        "passphrase": "",
        "id": "",
        "description": "autoadd",
        "stapler-class": "com.cloudbees.jenkins.plugins.sshcredentials.impl.BasicSSHUserPrivateKey",
        "$class": "com.cloudbees.jenkins.plugins.sshcredentials.impl.BasicSSHUserPrivateKey"}
        }


    data = json.dumps(postData)

    req = urllib2.Request(post_url)
    urllib2.urlopen(req,urllib.urlencode({'json':data}))

if __name__ == "__main__":
    jenkins_home=os.getenv('JENKINS_HOME')
    nodeIp = os.getenv('jenkins_node_ip')
    nodeSSHPort = os.getenv('jenkins_node_ssh_port')
    jenkinsIp = os.getenv('jenkins_host_ip')
    jenkinsPort = os.getenv('jenkins_host_port')

    if (nodeIp is None) or (nodeSSHPort is None) or (jenkinsIp is None) or (jenkinsPort is None):
        try:
            print ("env is None: jenkins_node_ip, jenkins_node_port, jenkins_host_ip, jenkins_host_port")
            os._exit(0)
        except:
            print ("Failed to exit the proccess")
    
    print ("variables........................")
    print ("nodeIp	: %s" %nodeIp)
    print ("nodeSSHPort	: %s" %nodeSSHPort)
    print ("jenkinsIp	: %s" %jenkinsIp)
    print ("jeninsPort	: %s" %jenkinsPort)
    print ("jenkins_home: %s" %jenkins_home)

    jenkinsMaster = "http://" + jenkinsIp + ":" + jenkinsPort
    while ((os.path.exists("%s/credentials.xml" %jenkins_home))== False):
        time.sleep(2)
        curl(jenkinsMaster)
    cid = getCredentialId(jenkins_home)
    server = jenkins.Jenkins(jenkinsMaster, username='', password="")
    version = server.get_version()

    print (version)

    params = {
        'port': nodeSSHPort,
        'username': 'juser',
        'credentialsId': cid,
        'host': nodeIp
    }

    server.create_node(
        'slave',
        nodeDescription='add slave',
        remoteFS='/var/jenkins',
        labels='precise',
        exclusive=True,
        launcher=jenkins.LAUNCHER_SSH,
        launcher_params=params)
