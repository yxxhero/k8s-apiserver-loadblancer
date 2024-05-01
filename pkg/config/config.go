package config

import (
	"fmt"

	"github.com/yxxhero/k8s-apiserver-loadblancer/pkg/util"
)

// Config represents the configuration of the application
type Config struct {
	ServiceName      string
	ServiceNamespace string
	ServiceType      string
	Kubeconfig       string
	StopCh           chan struct{}
	ID               string
}

// NewConfig returns a new Config
func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Verify() error {
	if !util.VerifyServiceType(c.ServiceType) {
		return fmt.Errorf("invalid service type: %s", c.ServiceType)
	}
	if c.ServiceName == "" {
		return fmt.Errorf("service name is required")
	}
	if c.ServiceNamespace == "" {
		return fmt.Errorf("service namespace is required")
	}
	if c.ID == "" {
		return fmt.Errorf("id is required")
	}
	return nil
}
