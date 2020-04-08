package standalone

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	p := NewStandalone("./testconfig.yaml")

	config, err := p.LoadConfig()

	assert.True(t, len(config) == 1)

	assert.Nil(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "eventstore", config[0].Kind)
	assert.Equal(t, "myeventstore", config[0].Metadata.Name)
	assert.Equal(t, "eventstore.azure.tablestorage", config[0].Spec.Type)
	assert.Equal(t, "storageAccountName", config[0].Spec.Metadata[0].Name)
	assert.Equal(t, "testaccount", config[0].Spec.Metadata[0].Value)
	assert.Equal(t, "storageAccountKey", config[0].Spec.Metadata[1].Name)
	assert.Equal(t, "testaccountkey", config[0].Spec.Metadata[1].Value)
}
