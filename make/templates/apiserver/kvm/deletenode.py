import sys
import jenkins

def delete_jenkins_node(jenkinsmasterurl, nodename):
    #cid = '3ea2a0cd-b611-433b-9998-d2d8882f97b2'
    server = jenkins.Jenkins(jenkinsmasterurl)
    try:
        server.delete_node(nodename)
    except:
        print ('can not delete node')

def main():
    jenkinsmasterurl = sys.argv[1]
    nodename = sys.argv[2]
    delete_jenkins_node(jenkinsmasterurl, nodename) 

if __name__ == "__main__":
    main()
