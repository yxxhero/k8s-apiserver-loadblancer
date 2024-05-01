package k8s

import (
	"context"
	"encoding/json"

	v1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

// CreateService creates a service in the k8s cluster
func CreateService(clientset *kubernetes.Clientset, service *v1.Service) (*v1.Service, error) {
	result, err := clientset.CoreV1().Services(service.Namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetService gets a service in the k8s cluster
func GetService(clientset *kubernetes.Clientset, namespace, name string) (*v1.Service, error) {
	result, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// IsExistService checks if a service exists in the k8s cluster
func IsExistService(clientset *kubernetes.Clientset, namespace, name string) (bool, error) {
	_, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// SetServiceType sets the type of a service in the k8s cluster
func SetServiceType(clientset *kubernetes.Clientset, namespace, name, serviceType string) (bool, error) {
	service, err := GetService(clientset, namespace, name)
	if err != nil {
		return false, err
	}

	if service.Spec.Type == v1.ServiceType(serviceType) {
		return false, nil
	}

	patchData := map[string]interface{}{"spec": map[string]interface{}{"type": serviceType}}

	playloadBytes, err := json.Marshal(patchData)
	if err != nil {
		return false, err
	}

	_, err = clientset.CoreV1().Services(namespace).Patch(context.TODO(), name, types.MergePatchType, playloadBytes, metav1.PatchOptions{})

	if err != nil {
		return false, err
	}
	return true, nil
}
