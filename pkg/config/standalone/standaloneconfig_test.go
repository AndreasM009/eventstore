package standalone

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	p := NewStandalone("./testconfig.yaml")

	config, err := p.LoadConfig()
	assert.Nil(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "eventstore", config.Kind)
	assert.Equal(t, "myeventstore", config.Metadata.Name)
	assert.Equal(t, "eventstore.azure.tablestorage", config.Spec.Type)
	assert.Equal(t, "storageAccount", config.Spec.Metadata[0].Name)
	assert.Equal(t, "testaccount", config.Spec.Metadata[0].Value)
	assert.Equal(t, "storageAccountKey", config.Spec.Metadata[1].Name)
	assert.Equal(t, "testaccountkey", config.Spec.Metadata[1].Value)
}
