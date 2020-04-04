package operator

import (
	"context"
	"errors"
	"log"

	eventstore "github.com/AndreasM009/eventstore-service-go/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const (
	eventStoreWorker = "EventStore"
)

// Operator interface
type Operator interface {
	Run(context.Context) (<-chan struct{}, error)
}

type operator struct {
	eventstoreClient    *eventstore.Clientset
	kubernetesClient    *kubernetes.Clientset
	eventstoreInformer  cache.SharedIndexInformer
	eventstoreQueue     workqueue.RateLimitingInterface
	eventstoreWorker    QueueWorker
	eventstoreProcessor Processor
}

// NewOperator creates a new Eventstore Operator
func NewOperator(eventstoreClient *eventstore.Clientset, kubernetesClient *kubernetes.Clientset) Operator {
	op := &operator{
		kubernetesClient: kubernetesClient,
		eventstoreClient: eventstoreClient,
		eventstoreInformer: createEventstoreIndexInformer(
			context.TODO(), eventstoreClient, metav1.NamespaceAll, nil, nil),
		eventstoreQueue:     workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		eventstoreProcessor: newEventStoreProcessor(),
	}

	op.eventstoreWorker = newQueueWorker(
		eventStoreWorker, op.eventstoreInformer,
		op.eventstoreQueue, op.eventstoreProcessor.ProcessChanged, op.eventstoreProcessor.ProcessDeleted)

	op.eventstoreInformer.AddEventHandler(newInformerHandler(op.eventstoreQueue))

	return op
}

func (op *operator) Run(ctx context.Context) (<-chan struct{}, error) {
	stopContext, cancel := context.WithCancel(context.Background())

	go func() {
		defer utilruntime.HandleCrash()
		// stop worker
		defer op.eventstoreQueue.ShutDown()
		op.eventstoreInformer.Run(stopContext.Done())
		log.Println("Eventstore SharedIndexInformer stopped")
		cancel()
	}()

	if !cache.WaitForCacheSync(ctx.Done()) {
		err := errors.New("timed out waiting for caches to sync")
		utilruntime.HandleError(err)
		cancel()
		return stopContext.Done(), err
	}

	go func() {
		select {
		case <-ctx.Done():
			cancel()
			return
		case <-stopContext.Done():
			return
		}
	}()

	op.eventstoreWorker.Run(stopContext.Done())

	return stopContext.Done(), nil
}
