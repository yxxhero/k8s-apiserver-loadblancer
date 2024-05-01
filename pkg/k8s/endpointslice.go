package k8s

import (
	"context"

	"github.com/rs/zerolog/log"
	discoveryV1 "k8s.io/api/discovery/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CreateEndpointSliec creates a service in the k8s cluster
func CreateEndpointSliec(clientset *kubernetes.Clientset, endpointSlice *discoveryV1.EndpointSlice) (*discoveryV1.EndpointSlice, error) {
	result, err := clientset.DiscoveryV1().EndpointSlices(endpointSlice.Namespace).Create(context.TODO(), endpointSlice, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetEndpointSlice gets a service in the k8s cluster
func GetEndpointSlice(clientset *kubernetes.Clientset, namespace, name string) (*discoveryV1.EndpointSlice, error) {
	result, err := clientset.DiscoveryV1().EndpointSlices(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// IsExistEndpointSlice checks if a service exists in the k8s cluster
func IsExistEndpointSlice(clientset *kubernetes.Clientset, namespace, name string) (bool, error) {
	_, err := clientset.DiscoveryV1().EndpointSlices(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			log.Info().Msgf("EndpointSlice %s in namespace %s does not exist", name, namespace)
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// UpdateEndpointSlice updates a service in the k8s cluster
func UpdateEndpointSlice(clientset *kubernetes.Clientset, endpointSlice *discoveryV1.EndpointSlice) (*discoveryV1.EndpointSlice, error) {
	result, err := clientset.DiscoveryV1().EndpointSlices(endpointSlice.Namespace).Update(context.TODO(), endpointSlice, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	return result, nil
}
