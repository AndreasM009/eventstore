package operator

import (
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

func newInformerHandler(queue workqueue.RateLimitingInterface) cache.ResourceEventHandlerFuncs {
	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			queue.Add(obj)
		},
		UpdateFunc: func(_, obj interface{}) {
			queue.Add(obj)
		},
		DeleteFunc: func(obj interface{}) {
			queue.Add(obj)
		},
	}
}
