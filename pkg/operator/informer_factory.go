package operator

import (
	"context"

	eventstorev1alphav1 "github.com/AndreasM009/eventstore/pkg/apis/eventstore/v1alpha1"
	scheme "github.com/AndreasM009/eventstore/pkg/client/clientset/versioned"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
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

func createDeploymentIndexInformer(
	ctx context.Context,
	kubernetesClient kubernetes.Interface,
	namespace string,
	fieldSelector fields.Selector,
	labelSelector labels.Selector) cache.SharedIndexInformer {
	deploymentsClient := kubernetesClient.AppsV1().Deployments(namespace)

	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if fieldSelector != nil {
					options.FieldSelector = fieldSelector.String()
				}
				if labelSelector != nil {
					options.LabelSelector = labelSelector.String()
				}
				return deploymentsClient.List(ctx, options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if fieldSelector != nil {
					options.FieldSelector = fieldSelector.String()
				}
				if labelSelector != nil {
					options.LabelSelector = labelSelector.String()
				}
				return deploymentsClient.Watch(ctx, options)
			},
		},
		&appsv1.Deployment{},
		0,
		cache.Indexers{},
	)
}
