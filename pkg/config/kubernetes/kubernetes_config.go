package kubernetes

import (
	"errors"
	"strings"

	"github.com/AndreasM009/eventstore-service-go/pkg/config"
)

type kubernetesConfigurationProvider struct {
	// todo: add client to request Configuration from ControlPlane (operator)
	evenstoreNames []string
}

// NewKubernetes creates a new Kubernetes ConfigurationProvider
func NewKubernetes(eventstoreNames string) (config.ConfigurationProvider, error) {
	n := strings.Split(strings.Trim(eventstoreNames, "'"), ",")
	if n[0] == "" {
		return nil, errors.New("no evenstores defined")
	}

	names := make([]string, len(n))
	for i, s := range n {
		names[i] = strings.TrimSpace(s)
	}

	return &kubernetesConfigurationProvider{
		evenstoreNames: names,
	}, nil
}

func (k *kubernetesConfigurationProvider) LoadConfig() (*config.Configuration, error) {
	// todo: use client to request configuration
	config := &config.Configuration{
		Kind: "eventstore",
		Metadata: config.ConfigurationMetadata{
			Name: "myeventstore",
		},
		Spec: config.Spec{
			Type: "eventstore.azure.tablestorage",
			Metadata: []config.SpecMetadata{
				config.SpecMetadata{
					Name:  "storageAccountName",
					Value: "",
				},
				config.SpecMetadata{
					Name:  "storageAccountKey",
					Value: "",
				},
				config.SpecMetadata{
					Name:  "tableNameSuffix",
					Value: "",
				},
			},
		},
	}
	return config, nil
}
