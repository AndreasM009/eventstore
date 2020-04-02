package config

// ConfigurationMetadata metatdata props of config
type ConfigurationMetadata struct {
	Name string `yaml:"name"`
}

// SpecMetadata spec metadata part
type SpecMetadata struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// Spec spec part of config
type Spec struct {
	Type     string         `yaml:"type"`
	Metadata []SpecMetadata `yaml:"metadata"`
}

// Configuration for evenstore to use
type Configuration struct {
	Kind     string                `yaml:"kind"`
	Metadata ConfigurationMetadata `yaml:"metadata"`
	Spec     Spec                  `yaml:"spec"`
}

// ConfigurationProvider interface
type ConfigurationProvider interface {
	LoadConfig() (*Configuration, error)
}
