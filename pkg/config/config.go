package config

// Configuration for evenstore to use
type Configuration struct {
	Kind     string `yaml:"kind"`
	Metadata struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Type     string `yaml:"type"`
		Metadata []struct {
			Name  string `yaml:"name"`
			Value string `yaml:"value"`
		} `yaml:"metadata"`
	} `yaml:"spec"`
}

// ConfigurationProvider interface
type ConfigurationProvider interface {
	LoadConfig() (*Configuration, error)
}
