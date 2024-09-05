package mirror

import (
	"context"

	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/core/v1"
	discoveryV1 "k8s.io/api/discovery/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/yxxhero/k8s-apiserver-loadblancer/pkg/config"
	"github.com/yxxhero/k8s-apiserver-loadblancer/pkg/constant"
	"github.com/yxxhero/k8s-apiserver-loadblancer/pkg/informer"
	"github.com/yxxhero/k8s-apiserver-loadblancer/pkg/k8s"
)

// Run is the entry point for the mirror command
func Run(ctx context.Context, c *config.Config, k8sClient *kubernetes.Clientset) error {
	informer := informer.NewEndpointsliceInformer(constant.KubernetesServiceName, constant.KubeconfigServiceNamespace)

	if err := informer.Run(ctx, k8sClient, c.StopCh); err != nil {
		return err
	}

	processNextItem := func() bool {
		// Wait until there's a new item in the work queue
		key, quit := informer.WorkQueue.Get()
		if quit {
			return false
		}
		defer informer.WorkQueue.Done(key)
		es, err := informer.GetEndpointslice(key.(string))
		if err != nil {
			log.Error().Err(err).Msg("failed to get endpointSlice")
		}
		if es == nil {
			informer.WorkQueue.Forget(key)
			return true
		}
		err = reconcile(es, c, k8sClient)
		if err != nil {
			log.Error().Err(err).Msg("failed to reconcile")
			informer.WorkQueue.AddRateLimited(key)
			return true
		}

		return true
	}
	// Loop until the stop channel is closed
	for processNextItem() {
	}
	return nil
}

func convertSVCPorts(svcPorts []v1.ServicePort) []v1.ServicePort {
	result := []v1.ServicePort{}
	for _, port := range svcPorts {
		if port.Name == "https" && port.Port == 443 {
			result = append(result, v1.ServicePort{
				Name:       port.Name,
				Port:       6443,
				TargetPort: port.TargetPort,
				Protocol:   port.Protocol,
			})
			continue
		}
		result = append(result, port)
	}
	return result
}

// reconcile is the function that will be called to reconcile the state of the system
func reconcile(es *discoveryV1.EndpointSlice, c *config.Config, k8sclient *kubernetes.Clientset) error {
	ok, err := k8s.IsExistService(k8sclient, c.ServiceNamespace, c.ServiceName)
	if err != nil {
		log.Error().Err(err).Msgf("failed to check if service %s in namespace %s exists", c.ServiceName, c.ServiceNamespace)
		return err
	}

	s, err := k8s.GetService(k8sclient, constant.KubeconfigServiceNamespace, constant.KubernetesServiceName)
	if err != nil {
		log.Error().Err(err).Msgf("failed to get service %s in namespace %s", c.ServiceName, c.ServiceNamespace)
		return err
	}
	customSVC := s.DeepCopy()
	customSVC.Name = c.ServiceName
	customSVC.Namespace = c.ServiceNamespace
	customSVC.Spec.ClusterIPs = nil
	customSVC.Spec.ClusterIP = ""
	customSVC.UID = ""
	customSVC.ResourceVersion = ""
	customSVC.Spec.Ports = convertSVCPorts(s.Spec.Ports)
	customSVC.Spec.Type = v1.ServiceType(c.ServiceType)

	if !ok {
		_, err = k8s.CreateService(k8sclient, customSVC)
		if err != nil {
			log.Error().Err(err).Msgf("failed to create service %s in namespace %s", c.ServiceName, c.ServiceNamespace)
			return err
		}
	} else {
		_, err = k8s.SetServiceType(k8sclient, c.ServiceNamespace, c.ServiceName, c.ServiceType)
		if err != nil {
			log.Error().Err(err).Msgf("failed to update service %s in namespace %s", c.ServiceName, c.ServiceNamespace)
			return err
		}
	}

	ok, err = k8s.IsExistEndpointSlice(k8sclient, c.ServiceNamespace, c.ServiceName)
	if err != nil {
		log.Error().Err(err).Msgf("failed to check if endpointSlice %s in namespace %s exists", c.ServiceName, c.ServiceNamespace)
		return err
	}

	customES := es.DeepCopy()
	customES.Name = c.ServiceName
	customES.Namespace = c.ServiceNamespace
	customES.Labels["kubernetes.io/service-name"] = c.ServiceName
	customES.UID = ""
	customES.ResourceVersion = ""

	if !ok {
		_, err = k8s.CreateEndpointSliec(k8sclient, customES)
		if err != nil {
			log.Error().Err(err).Msgf("failed to create endpointSlice %s in namespace %s", c.ServiceName, c.ServiceNamespace)
			return err
		}
		return nil
	}
	_, err = k8s.UpdateEndpointSlice(k8sclient, customES)
	return err
}
