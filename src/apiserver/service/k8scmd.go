package service

import (
	"encoding/json"
	"errors"

	"github.com/inspursoft/board/src/common/utils"
	"io/ioutil"
	"net/http"

	//"k8s.io/client-go/kubernetes"
	//"k8s.io/client-go/pkg/api/resource"
	//"k8s.io/client-go/pkg/api/v1"
	//"k8s.io/client-go/rest"
	//apiCli "k8s.io/client-go/tools/clientcmd/api"

	"github.com/inspursoft/board/src/common/k8sassist"

	//"github.com/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"
	//"github.com/golang/glog"
	//"k8s.io/client-go/tools/clientcmd"
)

var (
	KubeMasterURL  = utils.GetConfig("KUBE_MASTER_URL")
	kubeConfigPath = utils.GetConfig("KUBE_CONFIG_PATH")
	EntryMethod    EntryMethodEnum
	CAPath         string
	TokenStr       string
)

type EntryMethodEnum int

const (
	Insecure EntryMethodEnum = iota + 1
	CertificateAuthority
	Token
)
const defaultEntry = Insecure

/* use k8sassit
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
*/

func Suspend(nodeName string) (bool, error) {

	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath(),
	})
	n := k8sclient.AppV1().Node()

	nodeData, err := n.Get(nodeName)
	nodeData.Unschedulable = true
	res, err := n.Update(nodeData)
	return res.Unschedulable, err

}

func Resume(nodeName string) (bool, error) {
	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath(),
	})
	n := k8sclient.AppV1().Node()
	nodeData, err := n.Get(nodeName)
	nodeData.Unschedulable = false
	res, err := n.Update(nodeData)
	return !res.Unschedulable, err
}

func Taint(nodeName string, effect string) error {
	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath(),
	})
	n := k8sclient.AppV1().Node()
	nodeData, err := n.Get(nodeName)
	logs.Info(nodeData)
	//TODO force drain all pods in this node
	return err
}

//get resource form k8s api-server
func k8sGet(resource interface{}, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusNotFound {
		return errors.New("StatusNotFound:" + string(body))
	}
	err = json.Unmarshal(body, resource)
	if err != nil {
		return err
	}
	return nil
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

/* Not support PV in this time
// setNFSVol is the function of setting a PVC to bound PV storage with nfs.
// The name is name of PVC, path is the path of nfs, cap is the capacity of PV and PVC
func SetNFSVol(name string, server, path string, cap int64) error {
	// common date pv and pvc
	var (
		// initialize date map
		storage v1.ResourceList = make(v1.ResourceList)
		// q is capacity of storage
		q resource.Quantity
		// storage mode. It can be set "ReadWriteOnce", "ReadOnlyMany" and "ReadWriteMany"
		mode []v1.PersistentVolumeAccessMode = []v1.PersistentVolumeAccessMode{v1.ReadWriteMany}
	)
	// init common date
	q.Set(cap)
	storage[v1.ResourceStorage] = q

	// get k8s client
	cli, err := K8sCliFactory("", KubeMasterURL(), "v1")
	apiSet, err := kubernetes.NewForConfig(cli)
	if err != nil {
		return err
	}
	// bound k8s client and create source
	pvSet := apiSet.PersistentVolumes()

	pv := v1.PersistentVolume{}

	pv.Name = name

	pv.Spec.NFS = &v1.NFSVolumeSource{
		Server:   server,
		Path:     path,
		ReadOnly: false,
	}
	pv.Spec.AccessModes = mode

	pv.Spec.Capacity = storage

	info, err := pvSet.Create(&pv)
	// set info logs with creating pv
	glog.Infof("%s", info)
	if err != nil {
		return err
	}

	pvcSet := apiSet.PersistentVolumeClaims("default")
	pvc := v1.PersistentVolumeClaim{}

	pvc.Name = name

	pvc.Spec.AccessModes = mode

	pvc.Spec.Resources.Requests = storage
	// set info logs with creating pvc
	infoP, err := pvcSet.Create(&pvc)
	glog.Infof("%s", infoP)
	if err != nil {
		return err
	}
	return nil
}
*/
