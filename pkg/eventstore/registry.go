package eventstore

import (
	"fmt"

	"github.com/AndreasM009/eventstore-go/store/azure/tablestorage"

	"github.com/AndreasM009/eventstore-go/store/inmemory"

	"github.com/AndreasM009/eventstore-go/store"
)

// Registry interface
type Registry interface {
	Create(name string) (store.EventStore, error)
}

type eventstoreRegistry struct {
	factory map[string]func() store.EventStore
}

// NewRegistry creates a new registry
func NewRegistry() Registry {
	r := &eventstoreRegistry{
		factory: map[string]func() store.EventStore{},
	}

	r.factory["eventstore.inmemory"] = func() store.EventStore {
		return inmemory.NewStore()
	}

	r.factory["eventstore.azure.tablestorage"] = func() store.EventStore {
		return tablestorage.NewStore()
	}

	return r
}

func (r *eventstoreRegistry) Create(name string) (store.EventStore, error) {
	factory, ok := r.factory[name]

	if !ok {
		return nil, fmt.Errorf("registry: can't create eventstore %s", name)
	}

	return factory(), nil
}
