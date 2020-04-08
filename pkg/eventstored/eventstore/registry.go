package eventstore

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/AndreasM009/eventstore-go/store/azure/tablestorage"

	"github.com/AndreasM009/eventstore-go/store/inmemory"

	"github.com/AndreasM009/eventstore-go/store"
	"github.com/AndreasM009/eventstore-service-go/pkg/eventstored/config"
)

// Registry interface
type Registry interface {
	Create(name string) (store.EventStore, error)
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

	return r
}

func (r *eventstoreRegistry) Create(name string) (store.EventStore, error) {
	factory, ok := r.factory[name]

	if !ok {
		return nil, fmt.Errorf("registry: can't create eventstore %s", name)
	}

	return factory(), nil
}

func (r *eventstoreRegistry) CreateFromConfiguration(configs []config.Configuration) (map[string]store.EventStore, error) {
	builder := strings.Builder{}
	resultmap := map[string]store.EventStore{}

	for _, v := range configs {
		s, err := r.Create(v.Spec.Type)
		if err != nil {
			builder.WriteString(fmt.Sprintf("%s\n", err))
		} else {
			metadata := store.Metadata{
				Properties: map[string]string{},
			}

			for _, m := range v.Spec.Metadata {
				metadata.Properties[m.Name] = m.Value
			}

			if err := s.Init(metadata); err != nil {
				builder.WriteString(fmt.Sprintf("%s\n", err))
			}

			log.Printf("registry: Eventstore '%s' initialized\n", v.Metadata.Name)
			resultmap[v.Metadata.Name] = s
		}
	}

	if builder.Len() != 0 {
		return resultmap, errors.New(builder.String())
	}

	return resultmap, nil
}
