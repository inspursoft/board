#!/usr/bin/python
 
import libvirt
import sys
 
def createConnection():
    conn = libvirt.openReadOnly(None)
    if conn == None:
        print 'Failed to open connection to QEMU/KVM'
        sys.exit(1)
    else:
        print '-----Connection is created successfully-----'
        return conn
 
def closeConnection(conn):
    print ''
    try:
        conn.close()
    except:
        print 'Failed to close the connection'
        return 1
    print 'Connection is closed'
 
def getDomInfoByName(conn, name):
    print ''
    print '----- get domain info by name -----'

    try:
        myDom = conn.lookupByName(name)
        return 0
    except:
        print 'Failed to find the domain with name "%s"' % name
        return 1
         
     
def getDomInfoByID(conn, id):
    print ''
    print '----- get domain info by ID -----'
    try:
        myDom = conn.lookupByID(id)
    except:
        print 'Failed to find the domain with ID "%d"' % id
        return 1
         
    print "Domain id is %d ; Name is %s" % (myDom.ID(), myDom.name())

def getDomInfoByName(conn, name):
    print ''
    print '----- get domain info by name -----'
    try:
        myDom = conn.lookupByName(name)
    except:
        print 'Failed to find the domain with name "%s"' % name
        return 1
    return myDom


