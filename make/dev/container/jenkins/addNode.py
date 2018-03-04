'''addNode.py $jenkins_home'''
import os,sys,re
import jenkins
import json
#import requests


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

def curl():

    username = "root"
    passwd = "123456a?"

    credential = {
    "scope": "GLOBAL",
    "username": "",
    "password": "",
    "description": "linda",
    "$class": "com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl"
  }
    credential["username"] = username
    credential["password"] = passwd

    credentialJson = (json.dumps(credential))

    cmd = '''curl -X POST 'http://@10.164.17.34:8085/credentials/store/system/domain/_/createCredentials' \
--data-urlencode 'json={
  "": "0",
  "credentials": %s
}'
''' %credentialJson
    print (cmd)
    os.system(cmd)


if __name__ == "__main__":
    pyPath = os.path.split(os.path.realpath(__file__))[0]
    jenkins_home=sys.argv[1]
    while ((os.path.exists("%s/credentials.xml" %jenkins_home))== False):
        curl()
    cid = getCredentialId(jenkins_home)


    server = jenkins.Jenkins('http://10.164.17.34:8085', username='admin', password="admin")
    version = server.get_version()
    print version

    params = {
        'port': '22',
        'username': 'juser',
        'credentialsId': cid,
        'host': '10.110.13.222'
    }
    server.create_node(
        'slave5',
        nodeDescription='my test slave',
        remoteFS='/home/juser',
        labels='precise',
        exclusive=True,
        launcher=jenkins.LAUNCHER_SSH,
        launcher_params=params)
