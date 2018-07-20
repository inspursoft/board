#!/bin/bash
python ../kvm/kvmnode.py "http://$jenkins_host_ip:$jenkins_host_port" $$1  $$2 "http://$jenkins_node_ip:$kvm_registry_port"
#curl "http://$jenkins_host_ip:$jenkins_host_port/job/$$1/buildWithParameters?node_name=$$2"