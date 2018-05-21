package types

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	emptyServiceNameErr      = errors.New("Service name is empty")
	portMaxErr               = errors.New("Nodeprot config exceeds the max limit")
	portMinErr               = errors.New("Nodeprot config exceeds the min limit")
	emptyDeployNameErr       = errors.New("Deployment name is empty")
	invalidReplicasErr       = errors.New("Replicas value is invalid")
	emptyContainerErr        = errors.New("Container config in yaml is empty")
	namespaceInconsistentErr = errors.New("Namespace value isn't consistent with project name.")
	deploymentKindErr        = errors.New("Deployment kind is invalid.")
	serviceKindErr           = errors.New("Service kind is invalid.")
	deploymentAPIVersionErr  = errors.New("Deployment API version is invalid.")
	serviceAPIVersionErr     = errors.New("Service API version is invalid.")
)

func CheckDeploymentConfig(projectName string, deployment Deployment) error {
	//check empty
	if deployment.Kind != deploymentKind {
		return deploymentKindErr
	}
	if deployment.APIVersion != deploymentAPIVersion {
		return deploymentAPIVersionErr
	}
	if deployment.ObjectMeta.Name == "" {
		return emptyDeployNameErr
	}
	if deployment.ObjectMeta.Namespace != projectName {
		return namespaceInconsistentErr
	}
	if *deployment.Spec.Replicas < 1 {
		return invalidReplicasErr
	}
	if len(deployment.Spec.Template.Spec.Containers) < 1 {
		return emptyContainerErr
	}

	for _, cont := range deployment.Spec.Template.Spec.Containers {

		err := checkStringHasUpper(cont.Name, cont.Image)
		if err != nil {
			return err
		}

		for _, com := range cont.Command {
			err := checkStringHasUpper(com)
			if err != nil {
				return err
			}
		}

		for _, volMount := range cont.VolumeMounts {
			err := checkStringHasUpper(volMount.Name, volMount.MountPath)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func checkStringHasUpper(str ...string) error {
	for _, node := range str {
		isMatch, err := regexp.MatchString("[A-Z]", node)
		if err != nil {
			return err
		}
		if isMatch {
			errString := fmt.Sprintf(`string "%s" has upper charactor`, node)
			return errors.New(errString)
		}
	}
	return nil
}

//check parameter of service yaml file
func CheckServiceConfig(projectName string, serviceConfig Service) error {
	//check empty
	if serviceConfig.Kind != serviceKind {
		return deploymentKindErr
	}
	if serviceConfig.APIVersion != serviceAPIVersion {
		return deploymentAPIVersionErr
	}
	if serviceConfig.ObjectMeta.Name == "" {
		return emptyServiceNameErr
	}
	if serviceConfig.ObjectMeta.Namespace != projectName {
		return namespaceInconsistentErr
	}

	for _, external := range serviceConfig.Spec.Ports {
		if external.NodePort > maxPort {
			return portMaxErr
		} else if external.NodePort < minPort {
			return portMinErr
		}
	}

	return nil
}
