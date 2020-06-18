package eventstore

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/AndreasM009/eventstore-impl/store/azure/cosmosdb"

	"github.com/AndreasM009/eventstore-impl/store/azure/tablestorage"

	"github.com/AndreasM009/eventstore-impl/store/inmemory"

	"github.com/AndreasM009/eventstore-impl/store"
	"github.com/AndreasM009/eventstore/pkg/eventstored/config"
)

// Registry interface
type Registry interface {
	Create(cfg config.Configuration) (store.EventStore, error)
	CreateFromConfiguration(configs []config.Configuration) (map[string]store.EventStore, error)
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

	r.factory["eventstore.azure.cosmosdb"] = func() store.EventStore {
		return cosmosdb.NewStore()
	}

	return r
}

func (r *eventstoreRegistry) Create(cfg config.Configuration) (store.EventStore, error) {
	factory, ok := r.factory[cfg.Spec.Type]

	if !ok {
		return nil, fmt.Errorf("registry: can't create eventstore %s", cfg.Spec.Type)
	}

	s := factory()

	metadata := store.Metadata{
		Properties: map[string]string{},
	}

	for _, m := range cfg.Spec.Metadata {
		metadata.Properties[m.Name] = m.Value
	}

	if err := s.Init(metadata); err != nil {
		return s, err
	}

	return s, nil
}

func (r *eventstoreRegistry) CreateFromConfiguration(configs []config.Configuration) (map[string]store.EventStore, error) {
	builder := strings.Builder{}
	resultmap := map[string]store.EventStore{}

	for _, v := range configs {
		s, err := r.Create(v)
		if err != nil {
			builder.WriteString(fmt.Sprintf("%s\n", err))
			resultmap[v.Metadata.Name] = s
		} else {
			log.Printf("registry: Eventstore '%s' initialized\n", v.Metadata.Name)
			resultmap[v.Metadata.Name] = s
		}
	}

	if builder.Len() != 0 {
		return resultmap, errors.New(builder.String())
	}

	return resultmap, nil
}
