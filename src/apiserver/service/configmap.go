package service

import (
	//	"errors"
	//"fmt"
	//"github.com/inspursoft/board/src/common/dao"
	"github.com/inspursoft/board/src/common/k8sassist"
	"github.com/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"github.com/inspursoft/board/src/common/model"

	//"strings"

	"github.com/astaxie/beego/logs"
)

func GetConfigMapListByProject(projectname string) ([]*model.ConfigMapStruct, error) {

	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath(),
	})
	cmlist, err := k8sclient.AppV1().ConfigMap(projectname).List()
	if err != nil {
		logs.Debug("Failed to get configmap list %s", projectname)
		return nil, err
	}
	configmapList := make([]*model.ConfigMapStruct, 0)
	for _, cmk8s := range cmlist.Items {
		configmapList = append(configmapList, &model.ConfigMapStruct{Name: cmk8s.Name,
			Namespace: cmk8s.Namespace,
			DataList:  cmk8s.Data,
		})
	}

	return configmapList, nil
}

func GetConfigMapListByUser(userID int64) ([]*model.ConfigMapStruct, error) {
	cms := make([]*model.ConfigMapStruct, 0)
	query := model.Project{}
	projectList, err := GetProjectsByMember(query, userID)
	if err != nil {
		logs.Debug("Failed to get projects %d %v", userID, err)
		return nil, err
	}

	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath(),
	})

	for _, p := range projectList {
		cmlist, err := k8sclient.AppV1().ConfigMap(p.Name).List()
		if err != nil {
			logs.Debug("Failed to get configmap list %s", p.Name)
			return nil, err
		}
		configmapList := make([]*model.ConfigMapStruct, 0)
		for _, cmk8s := range cmlist.Items {
			configmapList = append(configmapList, &model.ConfigMapStruct{Name: cmk8s.Name,
				Namespace: cmk8s.Namespace,
				DataList:  cmk8s.Data,
			})
		}
		cms = append(cms, configmapList...)
	}
	return cms, nil
}

func CreateConfigMapK8s(cm *model.ConfigMapStruct) (*model.ConfigMap, error) {
	// add the configmap to k8s
	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath(),
	})
	var err error
	var cmk8s model.ConfigMap
	cmk8s.Name = cm.Name
	cmk8s.Namespace = cm.Namespace
	cmk8s.Data = cm.DataList
	//TODO add a default label to it
	newcm, err := k8sclient.AppV1().ConfigMap(cm.Namespace).Create(&cmk8s)
	if err != nil {
		logs.Debug("Failed to add ConfigMap to K8s %v %v", cmk8s, err)
		return nil, err
	}
	logs.Debug("Added ConfigMap to K8s")
	return newcm, nil
}

func DeleteConfigMapK8s(configname string, projectname string) error {
	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath(),
	})
	err := k8sclient.AppV1().ConfigMap(projectname).Delete(configname)
	if err != nil {
		if types.IsNotFoundError(err) {
			logs.Debug("Not found ConfigMap %s", configname)
		} else {
			return err
		}
	}
	return nil
}

func GetConfigMapK8s(configname string, projectname string) (*model.ConfigMap, error) {
	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath(),
	})
	cm, err := k8sclient.AppV1().ConfigMap(projectname).Get(configname)
	if err != nil {
		if types.IsNotFoundError(err) {
			logs.Debug("Not found ConfigMap %s", configname)
			return nil, nil
		} else {
			return nil, err
		}
	}
	return cm, nil
}

func UpdateConfigMapK8s(cm *model.ConfigMapStruct) (*model.ConfigMap, error) {
	// add the configmap to k8s
	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath(),
	})
	var err error
	var cmk8s model.ConfigMap
	cmk8s.Name = cm.Name
	cmk8s.Namespace = cm.Namespace
	cmk8s.Data = cm.DataList

	newcm, err := k8sclient.AppV1().ConfigMap(cm.Namespace).Update(&cmk8s)
	if err != nil {
		logs.Debug("Failed to update ConfigMap to K8s %v %v", cmk8s, err)
		return nil, err
	}
	logs.Debug("Updated ConfigMap to K8s")
	return newcm, nil
}
