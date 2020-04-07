package kubernetes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/AndreasM009/eventstore-service-go/pkg/eventstored/config"
)

type kubernetesConfigurationProvider struct {
	// todo: add client to request Configuration from ControlPlane (operator)
	evenstoreNames   []string
	operatorEndpoint string
}

// NewKubernetes creates a new Kubernetes ConfigurationProvider
func NewKubernetes(eventstoreNames, operatorEndpoint string) (config.ConfigurationProvider, error) {
	n := strings.Split(strings.Trim(eventstoreNames, "'"), ",")
	if n[0] == "" {
		return nil, errors.New("no evenstores defined")
	}

	names := make([]string, len(n))
	for i, s := range n {
		names[i] = strings.TrimSpace(s)
	}

	return &kubernetesConfigurationProvider{
		evenstoreNames:   names,
		operatorEndpoint: operatorEndpoint,
	}, nil
}

func (k *kubernetesConfigurationProvider) LoadConfig() (*config.Configuration, error) {
	url := fmt.Sprintf("%s/eventstores", k.operatorEndpoint)

	resp, err := http.Get(url)
	if err != nil {
		err = fmt.Errorf("Can't load configuration from %s: %v", url, err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("Cant't read body from response: %v", err)
		return nil, err
	}

	configs := []config.Configuration{}
	err = json.Unmarshal(body, &configs)
	if err != nil {
		err = fmt.Errorf("Can't deserialize config from json: %v", err)
		return nil, err
	}

	if len(configs) == 0 {
		return nil, errors.New("No configration available")
	}

	log.Println(string(body))

	data, err := json.Marshal(configs)
	if err != nil {
		log.Println("Error serilaize configs")
		return nil, nil
	}

	log.Println(string(data))

	return &configs[0], nil
}
