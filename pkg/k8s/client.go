package k8s

import (
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NewClient returns a Kubernetes clientset based on the provided kubeconfig.
// If the kubeconfig is empty, it uses the in-cluster configuration.
// Otherwise, it builds the configuration from the provided kubeconfig file.
// It returns the clientset and any error encountered during the process.
func NewClient(kubeconfig string) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	if kubeconfig == "" {
		log.Info().Msg("kubeconfig is empty, use InClusterConfig")
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}
