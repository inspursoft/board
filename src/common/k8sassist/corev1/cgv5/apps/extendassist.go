package apps

import (
	"github.com/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"

	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type extendserver struct {
	clientset *types.Clientset
}

func (e *extendserver) ListSelectRelatePods(infos []*model.K8sInfo) (*model.PodList, error) {
	if infos == nil {
		return nil, nil
	}
	// remove duplicates pods.
	var objPods = make(map[string]model.Pod)
	for _, info := range infos {
		pods, err := e.listOneWorkLoadRelatePods(info)
		if err != nil {
			return nil, err
		}
		if pods != nil && len(pods.Items) > 0 {
			for _, p := range pods.Items {
				objPods[p.Namespace+"/"+p.Name] = p
			}
		}
	}
	pods := make([]model.Pod, len(objPods))
	i := 0
	for _, v := range objPods {
		pods[i] = v
		i++
	}
	return &model.PodList{Items: pods}, nil
}

func (e *extendserver) listOneWorkLoadRelatePods(info *model.K8sInfo) (*model.PodList, error) {
	// if the info object is pod, so we return the pod directly
	var podList *v1.PodList
	var err error
	if info.Kind == "Pod" {
		var pod *v1.Pod
		pod, err = e.clientset.CoreV1().Pods(info.Namespace).Get(info.Name, metav1.GetOptions{})
		if pod != nil {
			podList = &v1.PodList{
				Items: []v1.Pod{*pod},
			}
		}
	} else {
		sel, find, serr := e.getSelectorFromObject(info)
		if serr != nil || !find {
			return nil, serr
		}
		if sel == nil {
			logs.Warn("the kubernetes %+v has no selector, so ignore it", info)
			return nil, nil
		}
		logs.Debug("the selector of k8s object %s/%s is %+v", info.Kind, info.Name, sel)
		opts := model.ListOptions{
			LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{MatchLabels: sel}),
		}
		podList, err = e.clientset.CoreV1().Pods(info.Namespace).List(types.ToK8sListOptions(opts))
	}
	if err != nil {
		logs.Error("list k8s object %+v relate pods failed. Err:%+v", info, err)
		return nil, err
	}

	modelPodList := types.FromK8sPodList(podList)
	return modelPodList, nil
}

func (e *extendserver) getSelectorFromObject(info *model.K8sInfo) (map[string]string, bool, error) {
	var path []string
	switch info.Kind {
	case "Deployment", "StatefulSet", "DaemonSet", "ReplicaSet":
		path = []string{"spec", "selector", "matchLabels"}
	case "Service", "ReplicationController":
		path = []string{"spec", "selector"}
	case "Job":
		//reget the object.
		job, err := e.clientset.BatchV1().Jobs(info.Namespace).Get(info.Name, metav1.GetOptions{})
		if err != nil {
			return nil, false, err
		}
		if job.Spec.Selector == nil {
			// no selector , so relate no pods
			return nil, false, nil
		}
		return job.Spec.Selector.MatchLabels, true, nil
	case "CronJob":
		// ignore cronjob.
	default:
		return nil, false, nil
	}

	return getJsonMapField(info.Source, path)
}

func getJsonMapField(source string, path []string) (map[string]string, bool, error) {
	object := map[string]interface{}{}
	err := json.Unmarshal([]byte(source), &object)
	if err != nil {
		return nil, false, err
	}
	selector, find, err := utils.GetNestedField(object, path...)
	if err != nil || !find {
		return nil, false, err
	}
	if sel, ok := selector.(map[string]interface{}); ok {
		selectormap := map[string]string{}
		for k, v := range sel {
			sv, ok := v.(string)
			if !ok {
				return nil, false, fmt.Errorf("the selector %T is not map[string]string", selector)
			}
			selectormap[k] = sv
		}
		return selectormap, true, nil
	}
	return nil, false, fmt.Errorf("the selector %T is not map[string]string", selector)
}

func NewExtend(clientset *types.Clientset) *extendserver {
	return &extendserver{
		clientset: clientset,
	}
}
