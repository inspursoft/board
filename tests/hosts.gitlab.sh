echo "127.0.0.1       myserver.com
127.0.0.1       smtp.myserver.com
127.0.0.1       mail.myserver.com
" >> /etc/hosts

# Change board config

pwd
cat make/config/apiserver/kubeconfig

echo "$KUBE_MASTER_IP"
echo "$KUBE_MASTER_PORT"

echo "cp make/config/apiserver/kubeconfig /tmp/"
pwd
cp make/config/apiserver/kubeconfig /tmp/config
cat /tmp/config
