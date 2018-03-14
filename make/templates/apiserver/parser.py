#!/bin/python
import os
import sys
import ConfigParser
import commands

def getParameters(cfgDir):
    cfgFile = os.path.join(cfgDir, 'META.cfg')
    conf = ConfigParser.ConfigParser()
    conf.read(cfgFile)
    options = conf.options("para")
    try:
        apiserver = conf.get('para','apiserver')
    except:
        apiserver = ""
    try:
        value = conf.get('para','value')
    except:
        value = ""
    try:
        flag = conf.get('para','flag').strip()
    except:
        flag = ""
    try:
        extra  = conf.get('para','extras')
    except:
        extra = ""
    try:
        user_id = conf.get('para', 'user_id')
    except:
        user_id = ""
    try:
        docker_registry = conf.get('para','docker_registry')
    except:
        docker_registry = ""
    file_name = conf.get('para','file_name')
    repoName = extra.split('/')[0]
#    file_name = os.path.join(value, file_name_tmp)

    print ("apiserver: %s" %apiserver)
    print ("value: %s" %value)
    print ("flag: %s" %flag)
    print ("extra: %s" %extra)
    print ("file_name: %s" %file_name)
    print ("docker_registry: %s " %docker_registry)

    processDir = os.path.join(cfgDir, repoName)
    print ('cfgDir %s' %cfgDir)
    os.chdir(cfgDir)
        
    return apiserver, value, flag, extra, file_name, docker_registry, user_id


def curlApiServer(apiserver, user_id, build_id):
    os.system("curl %s/api/v1/jenkins-job/%s/%s" %(apiserver, user_id, build_id))

def runProcessService(apiserver, value, extra, file_name):
    flag_0 = ''
    try:
        file_name_list = file_name.split(',')
        file_name_1 = file_name_list[0]
        file_name_2 = file_name_list[1]
    except:
        flag_0 = "error"
    
    try: 
        extra_list = extra.split(',')
        extra_1 = extra_list[0]
        extra_2 = extra_list[1]
    except:
        flag_0 = "error"

    if flag_0 == "error":
        print ("META.cfg file is not correct!!!")
    else: 
        print ('Process Server................') 
        print("curl -X POST -H 'Content-Type: application/yaml' --data-binary @%s/%s %s" %(value, file_name_1, extra_1))
        os.system("curl -X POST -H 'Content-Type: application/yaml' --data-binary @%s/%s %s" %(value, file_name_1, extra_1))
        print("curl -X POST -H 'Content-Type: application/yaml' --data-binary @%s/%s %s" %(value, file_name_2, extra_2))
        os.system("curl -X POST -H 'Content-Type: application/yaml' --data-binary @%s/%s %s" %(value, file_name_2, extra_2))
    
    #for i in range(1,2):
    #    real_extra = ("extra_%s" %i)
    #    real_file_name = ("file_name_%s" %i)
    #    print("curl -X POST -H 'Content-Type: application/yaml' --data-binary @%s/%s %s" %(value, real_file_name, real_extra))
    #    os.system("curl -X POST -H 'Content-Type: application/yaml' --data-binary @%s/%s %s" %(value, real_file_name, real_extra))
def runRollingUpdate(value, extra, file_name):    
    os.system("curl -X PATCH -H 'Content-Type: application/strategic-merge-patch+json' -d %s %s" %(value, extra))

def getCommandPath():
    status, output = commands.getstatusoutput('which docker')
    if status == 0:
        return output
    else:
        print ('not installed docker on this host!!!!!!!!!!')
def runProcessImage(extra, value, file_name, docker_registry, apiserver):
    docker = getCommandPath()
    print ('start to build process image')
    print ('file_name: %s' %(file_name))
    print ('extra: %s' %(extra))
    print ('docker_registry: %s' %(docker_registry))
    os.system("%s build -f %s -t %s ." %(docker, file_name, extra))
    os.system("%s tag %s %s/%s" %(docker, extra,docker_registry,extra))
    os.system("%s push %s/%s" %(docker, docker_registry, extra))
    os.system("%s rmi %s/%s" %(docker, docker_registry, extra))
    os.system("%s rmi %s" %(docker, extra))

def run(apiserver, value, flag, extra, file_name, docker_registry):
    path = os.path.split(os.path.realpath(__file__))[0]
    if flag.strip() == "process-image":
        file_name = os.path.join(value, file_name)
        runProcessImage(extra, value, file_name, docker_registry, apiserver)
    elif flag.strip() == "process-service":
        runProcessService(apiserver, value, extra, file_name)
#    elif flag.strip() == "rolling-update":
        #runRollingUpdate(value, extra, file_name)

   
def checkFile(cfgDir, build_id):
    cfgFile = os.path.join(cfgDir, 'META.cfg')
    if os.path.isfile(cfgFile):
       apiserver, value, flag, extra, file_name, docker_registry, user_id = getParameters(cfgDir)
       curlApiServer(apiserver, user_id, build_id)
       if flag in ['process-image', 'process-service']:
           run(apiserver, value, flag, extra, file_name, docker_registry)
    else:
       print ("the flag not in ['process-image', 'process-service']")
if __name__ == "__main__":
    base_repo_dir = sys.argv[1]
    workspace = sys.argv[2]
    build_id = sys.argv[3]
    repoName = sys.argv[4]
    metaDir=os.path.join(workspace, repoName)
    print ('meta file........%s' %metaDir)
    checkFile(metaDir, build_id)
