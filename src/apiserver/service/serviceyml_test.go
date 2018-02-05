package service

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/pkg/api/unversioned"
	modelK8s "k8s.io/client-go/pkg/api/v1"
	modelK8sExt "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

var serviceName = "service001"
var serviceConfig = Service{
	modelK8s.Service{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: modelK8s.ObjectMeta{
			Name:      serviceName,
			Namespace: "library",
		},
		Spec: modelK8s.ServiceSpec{
			Ports: []modelK8s.ServicePort{
				modelK8s.ServicePort{
					NodePort: int32(31080),
					Port:     int32(80),
				},
			},
			Selector: map[string]string{"app": serviceName},
			Type:     modelK8s.ServiceType("NodePort"),
		},
	},
}

var replicas int32 = 1
var image = "10.110.13.136:5000/library/mydemoshowing:1.0"
var deploymentConfig = Deployment{
	modelK8sExt.Deployment{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: modelK8s.ObjectMeta{
			Name:      serviceName,
			Namespace: "library",
		},
		Spec: modelK8sExt.DeploymentSpec{
			Replicas: &replicas,
			Template: modelK8s.PodTemplateSpec{
				ObjectMeta: modelK8s.ObjectMeta{
					Labels: map[string]string{"app": serviceName},
				},
				Spec: modelK8s.PodSpec{
					Containers: []modelK8s.Container{
						modelK8s.Container{
							Name:  "pod001",
							Image: image,
							Ports: []modelK8s.ContainerPort{
								modelK8s.ContainerPort{ContainerPort: 80},
							},
						},
					},
				},
			},
		},
	},
}

func TestCheckDeploymentPath(t *testing.T) {
	assert := assert.New(t)
	err := CheckDeploymentPath("./tmp")
	assert.Nil(err, "Error occurred while testing CheckDeploymentPath.")
	deleteFile("./tmp")
}

func TestServiceExists(t *testing.T) {
	assert := assert.New(t)
	s, _ := ServiceExists("", "")
	assert.False(s, "Error occurred while testing ServiceExists.")
}

func TestGenerateDeploymentYamlFileFromK8S(t *testing.T) {
	assert := assert.New(t)
	serviceURL := fmt.Sprintf("%s/apis/extensions/v1beta1/namespaces/%s/deployments/%s", kubeMasterURL(), deploymentConfig.Namespace, deploymentConfig.Name)
	err := GenerateDeploymentYamlFileFromK8S(serviceURL, "deployment.yaml")
	assert.Nil(err, "Error occurred while testing GenerateDeploymentYamlFileFromK8S.")
	err = DeleteServiceConfigYaml("deployment.yaml")
	assert.Nil(err, "Error occurred while testing DeleteServiceConfigYaml.")
}

func TestGenerateServiceYamlFileFromK8S(t *testing.T) {
	serviceURL := fmt.Sprintf("%s/api/v1/namespaces/%s/services/%s", kubeMasterURL(), serviceConfig.Namespace, serviceConfig.Name)
	assert := assert.New(t)
	err := GenerateServiceYamlFileFromK8S(serviceURL, "service.yaml")
	assert.Nil(err, "Error occurred while testing GenerateServiceYamlFileFromK8S.")
	err = DeleteServiceConfigYaml("service.yaml")
	assert.Nil(err, "Error occurred while testing DeleteServiceConfigYaml.")
}

func TestCheckServiceConfig(t *testing.T) {
	assert := assert.New(t)
	err := CheckServiceConfig("library", serviceConfig)
	assert.Nil(err, "Error occurred while testing CheckServiceConfig.")
}

func TestCheckDeploymentConfig(t *testing.T) {
	assert := assert.New(t)
	err := CheckDeploymentConfig("library", deploymentConfig)
	assert.Nil(err, "Error occurred while testing CheckDeploymentConfig.")
}
