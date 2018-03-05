'''addNode.py $jenkins_home'''
import os,sys,re
import jenkins
import json
import urllib
import urllib2


def getCredentialId(jenkins_home):
    credentialFile = os.path.join(jenkins_home,'credentials.xml')
    with open(credentialFile) as f:
        lines = f.readlines()
        for line in lines:
            m = (re.search('''<id>.*</id>''', line))
            if m:
                id_all = (re.search('''<id>.*</id>''', line).group(0))
                cid = (id_all.split('>')[1]).split('<')[0]
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
    response = urllib2.urlopen(req,urllib.urlencode({'json':data}))

    print response.read()


if __name__ == "__main__":
    pyPath = os.path.split(os.path.realpath(__file__))[0]
    jenkins_home=sys.argv[1]
    nodeIp = os.getenv('jenkins_node_ip')
    jenkinsIp = os.getenv('jenkins_host_ip')
    jenkinsPort = os.getenv('jenkins_host_port')
    jenkinsMaster = "http://" + jenkinsIp + ":" + jenkinsPort
    while ((os.path.exists("%s/credentials.xml" %jenkins_home))== False):
        curl(jenkinsMaster)
    cid = getCredentialId(jenkins_home)


    #server = jenkins.Jenkins('http://10.164.17.34:8085', username='admin', password="admin")
    server = jenkins.Jenkins(jenkinsMaster, username='admin', password="admin")
    version = server.get_version()
    print version

    params = {
        'port': '22',
        'username': 'juser',
        'credentialsId': cid,
        'host': ''
    }
    params["host"] = nodeIp

    server.create_node(
        'slave5',
        nodeDescription='my test slave',
        remoteFS='/home/juser',
        labels='precise',
        exclusive=True,
        launcher=jenkins.LAUNCHER_SSH,
        launcher_params=params)
