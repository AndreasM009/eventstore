package operator

import (
	"context"
	"errors"
	"log"

	eventstorev1alpha1 "github.com/AndreasM009/eventstore/pkg/apis/eventstore/v1alpha1"
	eventstore "github.com/AndreasM009/eventstore/pkg/client/clientset/versioned"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const (
	eventStoreWorker = "EventStore"
	deploymentWorker = "Deployment"
)

// Operator interface
type Operator interface {
	Run(context.Context) (<-chan struct{}, error)
	InitCustomResourceDefinitions() error
}

type operator struct {
	eventstoreClient    *eventstore.Clientset
	kubernetesClient    *kubernetes.Clientset
	extensionClient     *apiextensionsclient.Clientset
	eventstoreInformer  cache.SharedIndexInformer
	deploymentInformer  cache.SharedIndexInformer
	eventstoreQueue     workqueue.RateLimitingInterface
	deploymentQueue     workqueue.RateLimitingInterface
	eventstoreWorker    QueueWorker
	deploymentWorker    QueueWorker
	eventstoreProcessor Processor
	deploymentProcessor Processor
}

// NewOperator creates a new Eventstore Operator
func NewOperator(eventstoreClient *eventstore.Clientset, kubernetesClient *kubernetes.Clientset, extensionClient *apiextensionsclient.Clientset) Operator {
	op := &operator{
		kubernetesClient: kubernetesClient,
		eventstoreClient: eventstoreClient,
		extensionClient:  extensionClient,
		eventstoreInformer: createEventstoreIndexInformer(
			context.TODO(), eventstoreClient, metav1.NamespaceAll, nil, nil),
		eventstoreQueue:     workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		eventstoreProcessor: newEventStoreProcessor(kubernetesClient),
		deploymentInformer: createDeploymentIndexInformer(
			context.TODO(), kubernetesClient, metav1.NamespaceAll, nil, nil),
		deploymentQueue:     workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		deploymentProcessor: newDeploymentProcessor(kubernetesClient),
	}

	op.eventstoreWorker = newQueueWorker(
		eventStoreWorker, op.eventstoreInformer,
		op.eventstoreQueue, op.eventstoreProcessor.ProcessChanged, op.eventstoreProcessor.ProcessDeleted)

	op.eventstoreInformer.AddEventHandler(newInformerHandler(op.eventstoreQueue))

	op.deploymentWorker = newQueueWorker(
		deploymentWorker, op.deploymentInformer,
		op.deploymentQueue, op.deploymentProcessor.ProcessChanged, op.deploymentProcessor.ProcessDeleted)

	op.deploymentInformer.AddEventHandler(newInformerHandler(op.deploymentQueue))

	return op
}

func (op *operator) Run(ctx context.Context) (<-chan struct{}, error) {
	stopContext, cancel := context.WithCancel(context.Background())

	go func() {
		// stop worker
		defer op.eventstoreQueue.ShutDown()
		op.eventstoreInformer.Run(stopContext.Done())
		log.Println("Eventstore SharedIndexInformer stopped")
		cancel()
	}()

	go func() {
		// stop worker
		defer op.deploymentQueue.ShutDown()
		op.deploymentInformer.Run(stopContext.Done())
		log.Println("Deployment SharedIndexInformer stopped")
		cancel()
	}()

	if !cache.WaitForCacheSync(ctx.Done()) {
		err := errors.New("timed out waiting for caches to sync")
		utilruntime.HandleError(err)
		cancel()
		return stopContext.Done(), err
	}

	go func() {
		defer utilruntime.HandleCrash()

		select {
		case <-ctx.Done():
			cancel()
			return
		case <-stopContext.Done():
			return
		}
	}()

	op.eventstoreWorker.Run(stopContext.Done())
	op.deploymentWorker.Run(stopContext.Done())

	return stopContext.Done(), nil
}

// InitCustomResourceDefinitions create custom resources
func (op *operator) InitCustomResourceDefinitions() error {
	return eventstorev1alpha1.CreateCustomResourceDefinition("", op.extensionClient)
}
