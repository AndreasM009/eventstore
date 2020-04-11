package runtime

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/AndreasM009/eventstore-service-go/pkg/eventstored/config"

	kubernetesConfig "github.com/AndreasM009/eventstore-service-go/pkg/eventstored/config/kubernetes"
	standaloneConfig "github.com/AndreasM009/eventstore-service-go/pkg/eventstored/config/standalone"

	"github.com/AndreasM009/eventstore-go/store"
	"github.com/AndreasM009/eventstore-service-go/pkg/eventstored/eventstore"
	"github.com/AndreasM009/eventstore-service-go/pkg/eventstored/http"
)

const (
	modeKubernetes = "kubernetes"
	modeStandalone = "standalone"
)

var (
	modeFlag              = flag.String("mode", "standalone", "Run mode: 'standalone' or 'kubernetes'")
	portFlag              = flag.Int("port", 5000, "Server port to use")
	configFilePathFlag    = flag.String("config", "", "Path to config file (standalone only).")
	eventStoreNamesFlags  = flag.String("eventstores", "", "Comma separated names of eventstores that are associated with the Application Pod (Kubernetes only).")
	operatorEndpointFlags = flag.String("operatorendpoint", "", "Endpoint of operator control plane (kubernetes only).")
)

// Runtime interface to run an EventStore
type Runtime interface {
	FromFlags() error
	Start() error
}

type runtime struct {
	started  bool
	registry eventstore.Registry
	stores   map[string]store.EventStore
	server   http.Server
}

// NewRuntime creates a new EventStore runtime
func NewRuntime() Runtime {
	return &runtime{}
}

func (r *runtime) FromFlags() error {
	flag.Parse()

	var cfg []config.Configuration
	var err error

	switch *modeFlag {
	case modeStandalone:
		cfgProvider := standaloneConfig.NewStandalone(*configFilePathFlag)

		cfg, err = cfgProvider.LoadConfig()
		if err != nil {
			log.Println(err)
		}
	case modeKubernetes:
		cfgProvider, err := kubernetesConfig.NewKubernetes(*eventStoreNamesFlags, *operatorEndpointFlags)
		if err != nil {
			return err
		}
		cfg, err = cfgProvider.LoadConfig()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("runtime: unknown runtime mode %s", *modeFlag)
	}

	r.registry = eventstore.NewRegistry()
	r.stores, err = r.registry.CreateFromConfiguration(cfg)
	if err != nil {
		log.Printf("runtime: %s\n", err)
	}

	r.server = http.NewServer(*portFlag, r.stores, r.registry)

	return nil
}

func (r *runtime) Start() error {

	if r.server == nil || r.registry == nil {
		return errors.New("runtime: runtime not initialized correctly")
	}

	if r.started {
		return errors.New("runtime: runtime already started")
	}

	r.server.StartNonBlocking()
	r.started = true
	log.Printf("runtime: Started on port %v\n", *portFlag)
	return nil
}
