package standalone

import (
	"fmt"
	"io/ioutil"

	"github.com/AndreasM009/eventstore-service-go/pkg/eventstored/config"
	"gopkg.in/yaml.v2"
)

type standalongConfigurationProvider struct {
	configFilePath string
}

// NewStandalone creates a new standaloneConfigurationProvider
func NewStandalone(filePath string) config.ConfigurationProvider {
	return &standalongConfigurationProvider{
		configFilePath: filePath,
	}
}

func (p *standalongConfigurationProvider) LoadConfig() ([]config.Configuration, error) {
	data, err := ioutil.ReadFile(p.configFilePath)
	if err != nil {
		return nil, fmt.Errorf("standalone config: can't read config file: %s", err)
	}

	cnfg := config.Configuration{}

	err = yaml.Unmarshal(data, &cnfg)
	if err != nil {
		return nil, fmt.Errorf("standalone config: can't read yaml in config file: %s", err)
	}

	result := []config.Configuration{cnfg}

	return result, nil
}
