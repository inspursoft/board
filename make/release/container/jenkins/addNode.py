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

def createCredential(jenkinsMaster, nodeUsername, nodePasswd):
    para = "--data-urlencode"
    post_url = '%s/credentials/store/system/domain/_/createCredentials ' %jenkinsMaster
    credent = {
    "scope": "GLOBAL",
    "username": nodeUsername,
    "password": nodePasswd,
    "description": "auto added",
    "$class": "com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl"}

    data = json.dumps(credent)
    cmd = '''curl -X POST %s %s 'json={"":"0","credentials":%s}'
          ''' %(post_url, para, data)
    print (cmd)
    os.system(cmd)

if __name__ == "__main__":
    jenkins_home=os.getenv('JENKINS_HOME')
    nodeIp = os.getenv('jenkins_node_ip')
    nodeSSHPort = os.getenv('jenkins_node_ssh_port')
    jenkinsIp = os.getenv('jenkins_host_ip')
    jenkinsPort = os.getenv('jenkins_host_port')
    jenkinsNodeUsername = os.getenv('jenkins_node_username')
    jenkinsNodePasswd = os.getenv('jenkins_node_password')
    jenkinsNodeVolume = os.getenv('jenkins_node_volume')

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
    print ("jenkins_node_volume: %s" %jenkinsNodeVolume)

    jenkinsMaster = "http://" + jenkinsIp + ":" + jenkinsPort

    baseJob = jenkinsMaster + "/job/base"
    while (urllib.urlopen(baseJob).code>400):
        time.sleep(2)
   
    if ((os.path.exists("%s/credentials.xml" %jenkins_home))== False): 
        createCredential(jenkinsMaster, jenkinsNodeUsername, jenkinsNodePasswd)

    server = jenkins.Jenkins(jenkinsMaster, username="", password="")
    try:
        nodeinfo = server.get_node_info('slave')
    except:
        nodeinfo = ''
    if len(nodeinfo) > 0:
        os._exit(0)

    cid = getCredentialId(jenkins_home)
    version = server.get_version()

    print ("Jenkins verison is : %s" %version)

    params = {
        'port': nodeSSHPort,
        'username': 'juser',
        'credentialsId': cid,
        'host': nodeIp
    }
    try:
        server.create_node(
            'slave',
            nodeDescription='add slave',
            remoteFS=jenkinsNodeVolume,
            numExecutors=8,
            labels='slave',
            exclusive=False,
            launcher=jenkins.LAUNCHER_SSH,
            launcher_params=params)
    except Exception as e:
        print ("Create node failed", e)
