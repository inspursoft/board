package service

import (
	"encoding/json"
	"errors"

	"git/inspursoft/board/src/common/utils"
	"io/ioutil"
	"net/http"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	apiCli "k8s.io/client-go/tools/clientcmd/api"

	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeMasterURL             = utils.GetConfig("KUBE_MASTER_URL")
	EntryMethod   EntryMethodEnum
	CAPath        string
	TokenStr      string
)

type EntryMethodEnum int

const (
	Insecure EntryMethodEnum = iota + 1
	CertificateAuthority
	Token
)
const defaultEntry = Insecure

func K8sCliFactory(clusterName string, masterUrl string, apiVersion string) (*rest.Config, error) {
	cli := apiCli.NewConfig()
	cli.Clusters[clusterName] = &apiCli.Cluster{
		Server:     masterUrl,
		APIVersion: apiVersion}
	switch EntryMethod {
	case CertificateAuthority:
		cli.Clusters[clusterName].CertificateAuthority = CAPath
	case Token:
		cli.AuthInfos[clusterName].Token = TokenStr
	default:
		EntryMethod = defaultEntry
		cli.Clusters[clusterName].InsecureSkipTLSVerify = true
	}
	cli.CurrentContext = clusterName
	clientBuilder := clientcmd.NewNonInteractiveClientConfig(*cli, clusterName,
		&clientcmd.ConfigOverrides{}, nil)
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

//get resource form k8s api-server
func k8sGet(resource interface{}, url string) (bool, error) {
	resp, err := http.Get(url)

	if err != nil {
		return true, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return true, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return false, errors.New(string(body))
	}
	err = json.Unmarshal(body, resource)
	if err != nil {
		return true, err
	}

	return true, nil
}

//get resource form k8s api-server
func GetK8sData(resource interface{}, url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, resource)
	if err != nil {
		return body, err
	}

	return body, nil
}
