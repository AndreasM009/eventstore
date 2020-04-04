package operator

import (
	"context"

	eventstorev1alphav1 "github.com/AndreasM009/eventstore-service-go/pkg/apis/eventstore/v1alpha1"
	scheme "github.com/AndreasM009/eventstore-service-go/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

// createEventstoreIndexInformer creates a new SharedIndexInformer
func createEventstoreIndexInformer(
	ctx context.Context,
	eventstoreClient scheme.Interface,
	namespace string,
	fieldSelector fields.Selector,
	labelSelector labels.Selector) cache.SharedIndexInformer {
	evtClient := eventstoreClient.EventstoreV1alpha1().Eventstores(namespace)
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if fieldSelector != nil {
					options.FieldSelector = fieldSelector.String()
				}
				if labelSelector != nil {
					options.LabelSelector = labelSelector.String()
				}
				return evtClient.List(ctx, options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if fieldSelector != nil {
					options.FieldSelector = fieldSelector.String()
				}
				if labelSelector != nil {
					options.LabelSelector = labelSelector.String()
				}
				return evtClient.Watch(ctx, options)
			},
		},
		&eventstorev1alphav1.Eventstore{},
		0,
		cache.Indexers{},
	)
}
