package informer

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	discoveryV1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type EndpointliceInformer struct {
	selector  string
	namespace string
	WorkQueue workqueue.RateLimitingInterface
	indexer   cache.Indexer
}

func NewEndpointsliceInformer(name string, namespace string) *EndpointliceInformer {
	return &EndpointliceInformer{
		selector:  fmt.Sprintf("metadata.name=%s", name),
		WorkQueue: workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
	}
}

func (ei *EndpointliceInformer) GetEndpointslice(key string) (*discoveryV1.EndpointSlice, error) {
	ee, exists, err := ei.indexer.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}
	return ee.(*discoveryV1.EndpointSlice), nil
}

func (ei *EndpointliceInformer) Run(ctx context.Context, k8sClient kubernetes.Interface, stopChan chan struct{}) error {
	eeIndexer, eeInformer := cache.NewIndexerInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				options.FieldSelector = ei.selector

				return k8sClient.DiscoveryV1().EndpointSlices(ei.namespace).List(ctx, options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				options.FieldSelector = ei.selector
				return k8sClient.DiscoveryV1().EndpointSlices(ei.namespace).Watch(ctx, options)
			},
		},
		&discoveryV1.EndpointSlice{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				key, err := cache.MetaNamespaceKeyFunc(obj)
				if err != nil {
					log.Error().Err(err).Msg("Error getting key for object")
					return
				}
				ei.WorkQueue.Add(key)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				key, err := cache.MetaNamespaceKeyFunc(newObj)
				if err != nil {
					log.Error().Err(err).Msg("Error getting key for object")
				}
				ei.WorkQueue.Add(key)
			},
			DeleteFunc: func(obj interface{}) {
			},
		},
		cache.Indexers{},
	)

	ei.indexer = eeIndexer
	log.Info().Msg("Starting endpoint slice informer")
	go eeInformer.Run(stopChan)
	go func() {
		<-stopChan
		log.Info().Msg("Stopping endpoint slice informer")
		ei.WorkQueue.ShutDown()
	}()
	if !cache.WaitForCacheSync(stopChan, eeInformer.HasSynced) {
		return fmt.Errorf("timed out waiting for caches to sync")
	}
	return nil
}
