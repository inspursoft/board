#!/usr/bin/python
import os
import sys
import commands

def getNodeIp(vir):
    macline = commands.getoutput('virsh dumpxml %s|grep mac|grep address' %vir)
    mac = macline.split("'")[1]
    ipaddr = commands.getoutput('''arp -a|grep %s|cut -d '(' -f2|cut -d ')' -f1''' %mac)
    return ipaddr

if __name__ == "__main__":
   virs = sys.argv[1:9]
   for vir in virs:
       ipaddr = getNodeIp(vir)
       print (ipaddr)
