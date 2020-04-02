package injector

import (
	"github.com/kelseyhightower/envconfig"
)

// Config for eventstore sidecar injector
type Config struct {
	TLSCertFile            string `envconfig:"TLS_CERT_FILE" required:"true"`
	TLSKeyFile             string `envconfig:"TLS_KEY_FILE" required:"true"`
	SidecarImage           string `envconfig:"SIDECAR_IMAGE" required:"true"`
	SidecarImagePullPolicy string `envconfig:"SIDECAR_IMAGE_PULL_POLICY"`
	Namespace              string `envconfig:"NAMESPACE" required:"true"`
}

// NewConfig creates a default Config
func NewConfig() Config {
	return Config{
		SidecarImagePullPolicy: "Always",
	}
}

// NewConfigFromEnvironment loads Config from Environment Variables
func NewConfigFromEnvironment() (Config, error) {
	config := NewConfig()
	err := envconfig.Process("", &config)
	return config, err
}
