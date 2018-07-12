#!/bin/bash
python ../kvm/deletenode.py http://$jenkins_host_ip:$jenkins_host_port $$1 http://$jenkins_node_ip:$kvm_registry_port