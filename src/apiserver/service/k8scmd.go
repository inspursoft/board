package service

import (
	"git/inspursoft/board/src/common/utils"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	apiCli "k8s.io/client-go/tools/clientcmd/api"

	"k8s.io/client-go/tools/clientcmd"
)

var kubeMasterURL = utils.GetConfig("KUBE_MASTER_URL")

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

	cli, err := K8sCliFactory("", kubeMasterURL(), "v1")
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
	cli, err := K8sCliFactory("", kubeMasterURL(), "v1")
	apiSet, err := kubernetes.NewForConfig(cli)
	if err != nil {
		return false, err
	}
	n := apiSet.Nodes()
	nodeData, err := n.Get(nodeName)
	nodeData.Spec.Unschedulable = false
	res, err := n.Update(nodeData)
	return !res.Spec.Unschedulable, err
}
