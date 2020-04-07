package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
