package kubernetes

import (
	"fmt"
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
func TestEmptyEventstores(t *testing.T) {
	stores := ""

	cfg, err := NewKubernetes(stores, "")
	assert.NotNil(t, err)
	assert.Nil(t, cfg)
}

func TestSplitEventStores(t *testing.T) {
	stores := "a,b,c,d"

	cfg, err := NewKubernetes(stores, "")
	assert.Nil(t, err)
	assert.NotNil(t, cfg)

	prv := cfg.(*kubernetesConfigurationProvider)
	assert.Equal(t, 4, len(prv.evenstoreNames))
	assert.Equal(t, "a", prv.evenstoreNames[0])
	assert.Equal(t, "b", prv.evenstoreNames[1])
	assert.Equal(t, "c", prv.evenstoreNames[2])
	assert.Equal(t, "d", prv.evenstoreNames[3])
}

func TestSplitEventStoresSpaces(t *testing.T) {
	stores := "a, b, c, d"

	cfg, err := NewKubernetes(stores, "")
	assert.Nil(t, err)
	assert.NotNil(t, cfg)

	prv := cfg.(*kubernetesConfigurationProvider)
	assert.Equal(t, 4, len(prv.evenstoreNames))
	assert.Equal(t, "a", prv.evenstoreNames[0])
	assert.Equal(t, "b", prv.evenstoreNames[1])
	assert.Equal(t, "c", prv.evenstoreNames[2])
	assert.Equal(t, "d", prv.evenstoreNames[3])
}

func TestSplitEventStoresWithQuotes(t *testing.T) {
	stores := "'a,b,c,d'"

	cfg, err := NewKubernetes(stores, "")
	assert.Nil(t, err)
	assert.NotNil(t, cfg)

	prv := cfg.(*kubernetesConfigurationProvider)
	assert.Equal(t, 4, len(prv.evenstoreNames))
	assert.Equal(t, "a", prv.evenstoreNames[0])
	assert.Equal(t, "b", prv.evenstoreNames[1])
	assert.Equal(t, "c", prv.evenstoreNames[2])
	assert.Equal(t, "d", prv.evenstoreNames[3])
}

func TestFilterConfigs(t *testing.T) {
	stores := fmt.Sprintf("%s,%s", storeNameOne, storeNameTwo)

	data := createtestConfiguration()
	cfg, err := NewKubernetes(stores, "")

	assert.Nil(t, err)
	assert.NotNil(t, cfg)

	result := cfg.(*kubernetesConfigurationProvider).filterConfigs(data)

	assert.Equal(t, 2, len(*result))
}
