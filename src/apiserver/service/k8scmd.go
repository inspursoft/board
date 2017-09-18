package service

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	apiCli "k8s.io/client-go/tools/clientcmd/api"

	"fmt"

	"os"

	"k8s.io/client-go/tools/clientcmd"
)

var MasterUrl = fmt.Sprintf("%s:%s", os.Getenv("KUBE_IP"), os.Getenv("KUBE_PORT"))

func K8sCliFactory(clusterName string, masterUrl string, apiVersion string) (*rest.Config, error) {
	cli := apiCli.NewConfig()
	cli.Clusters[clusterName] = &apiCli.Cluster{
		Server:                masterUrl,
		InsecureSkipTLSVerify: true,
		APIVersion:            apiVersion}
	cli.CurrentContext = clusterName
	clientBuilder := clientcmd.NewNonInteractiveClientConfig(*cli, clusterName, &clientcmd.ConfigOverrides{}, nil)
	return clientBuilder.ClientConfig()
}
func Suspend(nodeName string) (bool, error) {

	cli, err := K8sCliFactory("", MasterUrl, "v1")
	apiSet, err := kubernetes.NewForConfig(cli)
	if err != nil {
		return false, err
	}
	n := apiSet.Nodes()
	nodeData, err := n.Get(nodeName)
	nodeData.Spec.Unschedulable = true
	res, err := n.Update(nodeData)
	return res.Spec.Unschedulable, err

}

func Resume(nodeName string) (bool, error) {

	cli, err := K8sCliFactory("", MasterUrl, "v1")
	apiSet, err := kubernetes.NewForConfig(cli)
	if err != nil {
		return false, err
	}
	n := apiSet.Nodes()
	nodeData, err := n.Get(nodeName)
	nodeData.Spec.Unschedulable = false
	res, err := n.Update(nodeData)
	if res.Spec.Unschedulable == false {
		return true, nil
	}
	return false, err

}
