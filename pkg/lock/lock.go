package lock

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection/resourcelock"

	"github.com/yxxhero/k8s-apiserver-loadblancer/pkg/config"
)

const (
	// LeaderElectionResourceLockName is the name of the resource lock
	LeaderElectionResourceLockName = "k8s-apiserver-leader-election"
)

func NewResourceLock(clientset *kubernetes.Clientset, config *config.Config) *resourcelock.LeaseLock {
	return &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      LeaderElectionResourceLockName,
			Namespace: config.ServiceNamespace,
		},
		Client: clientset.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: config.ID,
		},
	}
}
