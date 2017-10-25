package service

import (
	"testing"

	modelK8s "k8s.io/client-go/pkg/api/v1"
)

func Testk8sGet(t *testing.T) {
	var endpoint modelK8s.Endpoints
	_, err := k8sGet(&endpoint, endpointUrl)
	if err != nil {
		t.Errorf("Error occurred while test k8sGet: %+v\n", err)
	}
	if err == nil {
		t.Log("k8sGet is ok.\n")
	}

	_, err = k8sGet(&endpoint, noEndpointUrl)
	if err == nil {
		t.Errorf("Error occurred while test k8sGet: %+v\n", err)
	}
	if err != nil {
		t.Log("k8sGet is ok.\n")
	}

	_, err = k8sGet(&endpoint, invalidEndpointUrl)
	if err == nil {
		t.Errorf("Error occurred while test k8sGet: %+v\n", err)
	}
	if err != nil {
		t.Log("k8sGet is ok.\n")
	}
}

func TestGetK8sData(t *testing.T) {
	var endpoint modelK8s.Endpoints
	_, err := GetK8sData(&endpoint, endpointUrl)
	if err != nil {
		t.Errorf("Error occurred while test GetK8sData: %+v\n", err)
	}
	if err == nil {
		t.Log("GetK8sData is ok.\n")
	}

	_, err = GetK8sData(&endpoint, noEndpointUrl)
	if err != nil {
		t.Errorf("Error occurred while test GetK8sData: %+v\n", err)
	}
	if err == nil {
		t.Log("GetK8sData is ok.\n")
	}

	_, err = GetK8sData(&endpoint, invalidEndpointUrl)
	if err == nil {
		t.Errorf("Error occurred while test GetK8sData: %+v\n", err)
	}
	if err != nil {
		t.Log("GetK8sData is ok.\n")
	}
}
