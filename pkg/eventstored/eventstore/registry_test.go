package eventstore

import (
	"testing"

	"github.com/AndreasM009/eventstore-service-go/pkg/eventstored/config"
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
