#!/bin/bash
#python ../kvm/deletenode.py "http://$jenkins_host_ip:$jenkins_host_port" $$1 "http://$jenkins_node_ip:$kvm_registry_port"
curl "http://$jenkins_node_ip:$kvm_registry_port/release-node?job_name=$$1"
echo "Deleting node for job:$$1 at http://$jenkins_host_ip:$jenkins_host_port, KVM registry: http://$jenkins_node_ip:$kvm_registry_port"