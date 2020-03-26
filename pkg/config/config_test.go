package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

var testConfig = `
kind: eventstore
metadata:
  name: myeventstore
spec:
  type: eventstore.azure.tablestorage
  metadata:
  - name: storageAccount
    value: "testaccount"
  - name: storageAccountKey
    value: "testaccountkey"
`

func TestReadConfig(t *testing.T) {
	config := Configuration{}
	err := yaml.Unmarshal([]byte(testConfig), &config)
	assert.Nil(t, err)

	assert.Equal(t, "eventstore", config.Kind)
	assert.Equal(t, "myeventstore", config.Metadata.Name)
	assert.Equal(t, "eventstore.azure.tablestorage", config.Spec.Type)
	assert.Equal(t, "storageAccount", config.Spec.Metadata[0].Name)
	assert.Equal(t, "testaccount", config.Spec.Metadata[0].Value)
	assert.Equal(t, "storageAccountKey", config.Spec.Metadata[1].Name)
	assert.Equal(t, "testaccountkey", config.Spec.Metadata[1].Value)
}
