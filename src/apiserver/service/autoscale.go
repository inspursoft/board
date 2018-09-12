package service

import (
	"errors"
	"fmt"
	"git/inspursoft/board/src/common/dao"

	"git/inspursoft/board/src/common/k8sassist"
	"git/inspursoft/board/src/common/model"
)

func genAutoScaleObject(svc *model.ServiceStatus, autoscale *model.ServiceAutoScale) *model.AutoScale {
	minPod := int32(autoscale.MinPod)
	cpuPercent := int32(autoscale.CPUPercent)
	return &model.AutoScale{
		ObjectMeta: model.ObjectMeta{
			Name:      autoscale.HPAName,
			Namespace: svc.ProjectName,
		},
		Spec: model.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: model.CrossVersionObjectReference{
				Kind:       "Deployment",
				Name:       svc.Name,
				APIVersion: "apps/v1beta2",
			},
			MinReplicas:                    &minPod,
			MaxReplicas:                    int32(autoscale.MaxPod),
			TargetCPUUtilizationPercentage: &cpuPercent,
		},
	}
}

func CreateAutoScale(svc *model.ServiceStatus, autoscale *model.ServiceAutoScale) (*model.ServiceAutoScale, error) {
	// add the hpa from k8s
	hpa := genAutoScaleObject(svc, autoscale)
	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		K8sMasterURL: kubeMasterURL(),
	})
	var err error
	hpa, err = k8sclient.AppV1().AutoScale(hpa.Namespace).Create(hpa)
	if err != nil {
		return nil, err
	}

	// add the hpa from storage
	newhpa := *autoscale
	return &newhpa, nil
}

func ListAutoScales(svc *model.ServiceStatus) ([]model.ServiceAutoScale, error) {
	// list the hpas from storage
	return nil, nil
}

func UpdateAutoScale(svc *model.ServiceStatus, autoscale *model.ServiceAutoScale) (*model.ServiceAutoScale, error) {
	// update the hpa from k8s
	hpa := genAutoScaleObject(svc, autoscale)
	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		K8sMasterURL: kubeMasterURL(),
	})
	var err error
	hpa, err = k8sclient.AppV1().AutoScale(hpa.Namespace).Update(hpa)
	if err != nil {
		return nil, err
	}
	// update the hpa from storage
	newhpa := *autoscale
	return &newhpa, nil
}

func DeleteAutoScale(svc *model.ServiceStatus, hpaid int64) error {
	// delete the hpa from k8s
	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		K8sMasterURL: kubeMasterURL(),
	})
	err := k8sclient.AppV1().AutoScale(svc.ProjectName).Delete(fmt.Sprintf("%d", hpaid))
	if err != nil {
		return err
	}

	// delete the hpa from storage

	return nil
}

// AutoScale in database
func CreateAutoScaleDB(autoscale model.ServiceAutoScale) (int64, error) {
	autoscaleID, err := dao.AddAutoScale(autoscale)
	if err != nil {
		return 0, err
	}
	return autoscaleID, nil
}

func DeleteAutoScaleDB(autoscaleID int64) (bool, error) {
	s := model.ServiceAutoScale{ID: autoscaleID}
	_, err := dao.DeleteAutoScale(s)
	if err != nil {
		return false, err
	}
	return true, nil
}

func UpdateAutoScaleDB(autoscale model.ServiceAutoScale, fieldNames ...string) (bool, error) {
	if autoscale.ID == 0 {
		return false, errors.New("no AutoScale ID provided")
	}
	_, err := dao.UpdateAutoScale(autoscale, fieldNames...)
	if err != nil {
		return false, err
	}
	return true, nil
}
