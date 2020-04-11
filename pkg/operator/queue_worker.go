package operator

import (
	"fmt"
	"log"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

// QueueWorker interace
type QueueWorker interface {
	Run(stop <-chan struct{})
}

type changedEvent = func(obj interface{}) error
type deletedEvent = func(obj interface{}) error

type queueworker struct {
	name         string
	informer     cache.SharedIndexInformer
	queue        workqueue.RateLimitingInterface
	changedEvent changedEvent
	deletedEvent deletedEvent
}

func newQueueWorker(
	name string,
	informer cache.SharedIndexInformer,
	queue workqueue.RateLimitingInterface,
	changed changedEvent,
	deleted deletedEvent) QueueWorker {
	return &queueworker{
		name:         name,
		informer:     informer,
		queue:        queue,
		changedEvent: changed,
		deletedEvent: deleted,
	}
}

func (qw *queueworker) Run(stop <-chan struct{}) {
	go func() {
		wait.Until(qw.runWorker, time.Second, stop)
		log.Printf("Worker %s stopped\n", qw.name)
	}()
}

func (qw *queueworker) runWorker() {
	// loop until done
	for qw.processItems() {
	}
}

func (qw *queueworker) processItems() bool {
	obj, quit := qw.queue.Get()
	if quit {
		return false
	}

	// indicate queue that we have done a piece of work
	defer qw.queue.Done(obj)

	err := qw.processItem(obj)

	if err == nil {
		qw.queue.Forget(obj)
	} else if qw.queue.NumRequeues(obj) < 20 {
		qw.queue.AddRateLimited(obj)
	} else {
		qw.queue.Forget(obj)
		err := fmt.Errorf("Worker %s, error processing (giving up): %v", qw.name, err)
		log.Println(err)
		utilruntime.HandleError(err)
	}

	return true
}

func (qw *queueworker) processItem(obj interface{}) error {
	qobj, exists, err := qw.informer.GetIndexer().Get(obj)

	if err != nil {
		return fmt.Errorf("worker %s Error fetching object", qw.name)
	}

	if !exists {
		return qw.deletedEvent(obj)
	}

	return qw.changedEvent(qobj)
}
