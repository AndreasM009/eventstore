package operator

import "log"

type eventstoreProcessor struct {
}

func newEventStoreProcessor() Processor {
	return &eventstoreProcessor{}
}

func (p *eventstoreProcessor) ProcessChanged(obj interface{}) error {
	log.Println("Object changed")
	return nil
}

func (p *eventstoreProcessor) ProcessDeleted(obj interface{}) error {
	log.Println("Object deleted")
	return nil
}
