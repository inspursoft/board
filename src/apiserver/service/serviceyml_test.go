package service_test

var serviceName = "testservice001"

// var serviceConfig = Service{
// 	modelK8s.Service{
// 		TypeMeta: unversioned.TypeMeta{
// 			Kind:       "Service",
// 			APIVersion: "v1",
// 		},
// 		ObjectMeta: modelK8s.ObjectMeta{
// 			Name:      serviceName,
// 			Namespace: "library",
// 		},
// 		Spec: modelK8s.ServiceSpec{
// 			Ports: []modelK8s.ServicePort{
// 				modelK8s.ServicePort{
// 					NodePort: int32(31080),
// 					Port:     int32(80),
// 				},
// 			},
// 			Selector: map[string]string{"app": serviceName},
// 			Type:     modelK8s.ServiceType("NodePort"),
// 		},
// 	},
// }

var replicas int32 = 1
var image = "10.110.13.136:5000/library/mydemoshowing:1.0"

// var deploymentConfig = Deployment{
// 	modelK8sExt.Deployment{
// 		TypeMeta: unversioned.TypeMeta{
// 			Kind:       "Deployment",
// 			APIVersion: "extensions/v1beta1",
// 		},
// 		ObjectMeta: modelK8s.ObjectMeta{
// 			Name:      serviceName,
// 			Namespace: "library",
// 		},
// 		Spec: modelK8sExt.DeploymentSpec{
// 			Replicas: &replicas,
// 			Template: modelK8s.PodTemplateSpec{
// 				ObjectMeta: modelK8s.ObjectMeta{
// 					Labels: map[string]string{"app": serviceName},
// 				},
// 				Spec: modelK8s.PodSpec{
// 					Containers: []modelK8s.Container{
// 						modelK8s.Container{
// 							Name:  "pod001",
// 							Image: image,
// 							Ports: []modelK8s.ContainerPort{
// 								modelK8s.ContainerPort{ContainerPort: 80},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	},
// }
