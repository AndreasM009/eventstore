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
	key, quit := qw.queue.Get()
	if quit {
		return false
	}

	// indicate queue that we have done a piece of work
	defer qw.queue.Done(key)

	err := qw.processItem(key.(string))

	if err == nil {
		qw.queue.Forget(key)
	} else if qw.queue.NumRequeues(key) < 20 {
		qw.queue.AddRateLimited(key)
	} else {
		qw.queue.Forget(key)
		err := fmt.Errorf("Worker %s, error processing %s (giving up): %v", qw.name, key, err)
		log.Println(err)
		utilruntime.HandleError(err)
	}

	return true
}

func (qw *queueworker) processItem(key string) error {
	log.Printf("worker %s processing change to %s\n", qw.name, key)

	obj, exists, err := qw.informer.GetIndexer().GetByKey(key)

	if err != nil {
		return fmt.Errorf("worker %s Error fetching object with key %s", qw.name, key)
	}

	if !exists {
		log.Printf("Eventstore %s was deleted\n", key)
		return qw.deletedEvent(obj)
	}

	log.Printf("Eventstore %s was added or updated\n", key)
	return qw.changedEvent(obj)
}
