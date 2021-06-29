## Quick start with kind

About `kind` :

> kind is a tool for running local Kubernetes clusters using Docker container “nodes”.
> kind was primarily designed for testing Kubernetes itself, but may be used for local development or CI.

Install `kind` : https://kind.sigs.k8s.io/docs/user/quick-start/

Quick start Board with kind:

``` BASH
# prepare a three nodes cluster
$ cat > kind-config.yaml << EOF
# three node (two workers) cluster config
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  image: kindest/node:v1.18.4
- role: worker
  image: kindest/node:v1.18.4
- role: worker
  image: kindest/node:v1.18.4
networking:
  apiServerAddress: 127.0.0.1
EOF

# create the cluster
$ kind create cluster --name board --config kind-config.yaml
Creating cluster "board" ...
 ✓ Ensuring node image (kindest/node:v1.18.4) 🖼
 ✓ Preparing nodes 📦 📦 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
 ✓ Joining worker nodes 🚜
Set kubectl context to "kind-board"
You can now use your cluster with:

kubectl cluster-info --context kind-board

Have a question, bug, or feature request? Let us know! https://kind.sigs.k8s.io/#community 🙂

# If it works, you can get these
$ kubectl cluster-info --context kind-board
Kubernetes master is running at https://127.0.0.1:37939
KubeDNS is running at https://127.0.0.1:37939/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.

$ kubectl get no --context kind-board -owide
NAME                  STATUS   ROLES    AGE   VERSION   INTERNAL-IP   EXTERNAL-IP   OS-IMAGE           KERNEL-VERSION          CONTAINER-RUNTIME
board-control-plane   Ready    master   81s   v1.18.4   172.18.0.2    <none>        Ubuntu 20.04 LTS   3.10.0-693.el7.x86_64   containerd://1.4.0-beta.1-34-g49b0743c
board-worker          Ready    <none>   41s   v1.18.4   172.18.0.4    <none>        Ubuntu 20.04 LTS   3.10.0-693.el7.x86_64   containerd://1.4.0-beta.1-34-g49b0743c
board-worker2         Ready    <none>   41s   v1.18.4   172.18.0.3    <none>        Ubuntu 20.04 LTS   3.10.0-693.el7.x86_64   containerd://1.4.0-beta.1-34-g49b0743c

# prepare for Board cert
$ mkdir -p /etc/board/cert
$ docker cp board-control-plane:/etc/kubernetes/pki/ca.key /etc/board/cert/ca-key.pem
$ docker cp board-control-plane:/etc/kubernetes/pki/ca.crt /etc/board/cert/ca.pem

# prepare for components
$ mkdir board-k8s-requires
$ cd board-k8s-requires
# board-clusterrolebinding.yaml
$ wget https://raw.githubusercontent.com/inspursoft/board-installer/main/ansible_k8s/roles/kubectlCMD/templates/board-clusterrolebinding.yaml
$ kubectl apply -f board-clusterrolebinding.yaml --context kind-board

$ docker pull k8s.gcr.io/cadvisor:v0.30.2
$ kind load docker-image k8s.gcr.io/cadvisor:v0.30.2 --name board
# cadvisor.yaml
$ wget https://raw.githubusercontent.com/inspursoft/board-installer/main/ansible_k8s/roles/kubectlCMD/templates/cadvisor.yaml
# You need to replace the {{docker_dir}} placeholder in the file before deploying. The docker directory defaults to /var/lib/docker
$ kubectl apply -f cadvisor.yaml --context kind-board

# install Board
$ tar xvf board-offline-installer-VERSION[-ARCH].tgz -C /root
# edit Board config
$ cd /root/Deploy && vi board.cfg
# 
# hostname = 192.168.154.6        # set it to your host ip
# 
# kube_http_scheme = https
# kube_master_ip = 172.18.0.2     # set it as INTERNAL-IP of control-plane
# kube_master_port = 6443
# 
# registry_ip = 192.168.154.6     # set it to your registry ip
# registry_port = 5000            # set it to your registry port
# 
# devops_opt = legacy             # set it to legacy in test
# 
# gogits_host_ip = 192.168.154.6  # set it to your host ip
# 
# jenkins_host_ip = 192.168.154.6 # set it to your host ip
# 
# jenkins_node_password = 111111  # set it to your host password
 
# start
$ ./install.sh
# after install, connect the container to k8s
$ docker network connect kind archive_apiserver_1
# Visit Board on your browser.
# Wait for initialization, login to Board with admin/123456a?
# then you can have a preview on it!
```

## Uninstall

``` BASH
$ cd /root/Deploy
$ ./uninstall.sh

$ kind delete cluster --name=board
```
