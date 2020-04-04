package operator

// Processor interface to process CRD events
type Processor interface {
	ProcessChanged(obj interface{}) error
	ProcessDeleted(obj interface{}) error
}
