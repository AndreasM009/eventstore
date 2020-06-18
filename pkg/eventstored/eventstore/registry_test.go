package eventstore

import (
	"testing"

	"github.com/AndreasM009/eventstore/pkg/eventstored/config"
	"github.com/stretchr/testify/assert"
)

const (
	storeNameOne = "teststore-one"
	storeNameTwo = "teststore-two"
)

func createtestConfiguration() []config.Configuration {
	result := []config.Configuration{
		config.Configuration{
			Kind: "eventstore",
			Metadata: config.ConfigurationMetadata{
				Name: storeNameOne,
			},
			Spec: config.Spec{
				Type: "eventstore.inmemory",
				Metadata: []config.SpecMetadata{
					config.SpecMetadata{
						Name:  "Metadata-one",
						Value: "Metadata-value-one",
					},
					config.SpecMetadata{
						Name:  "Metadata-two",
						Value: "Metadata-value-two",
					},
				},
			},
		},
		config.Configuration{
			Kind: "eventstore",
			Metadata: config.ConfigurationMetadata{
				Name: storeNameTwo,
			},
			Spec: config.Spec{
				Type: "eventstore.inmemory",
				Metadata: []config.SpecMetadata{
					config.SpecMetadata{
						Name:  "Metadata-one",
						Value: "Metadata-value-one",
					},
					config.SpecMetadata{
						Name:  "Metadata-two",
						Value: "Metadata-value-two",
					},
				},
			},
		},
	}

	return result
}

func TestLoadConfiguration(t *testing.T) {
	registry := NewRegistry()
	testdata := createtestConfiguration()

	stores, err := registry.CreateFromConfiguration(testdata)

	assert.Nil(t, err)
	assert.NotNil(t, stores)
	assert.Equal(t, 2, len(stores))

	storeOne, ok := stores[storeNameOne]
	assert.True(t, ok)
	assert.NotNil(t, storeOne)

	storeTwo, ok := stores[storeNameTwo]
	assert.True(t, ok)
	assert.NotNil(t, storeTwo)
}

func TestLoadConfigurationWithBadConfig(t *testing.T) {
	registry := NewRegistry()
	testdata := createtestConfiguration()

	testdata = append(testdata, config.Configuration{
		Kind: "eventstore",
		Metadata: config.ConfigurationMetadata{
			Name: "NotImplemented",
		},
		Spec: config.Spec{
			Type: "eventstore.notimplementd",
			Metadata: []config.SpecMetadata{
				config.SpecMetadata{
					Name:  "Metadata-one",
					Value: "Metadata-value-one",
				},
				config.SpecMetadata{
					Name:  "Metadata-two",
					Value: "Metadata-value-two",
				},
			},
		},
	})

	stores, err := registry.CreateFromConfiguration(testdata)

	assert.NotNil(t, err)
	assert.NotNil(t, stores)
	assert.Equal(t, 3, len(stores))

	s, ok := stores["NotImplemented"]
	assert.Nil(t, s)
	assert.True(t, ok)
}
