package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Namespaces map[string]Deployment `yaml:"namespaces"`
}

type Deployment struct {
	Deployment map[string]DeploymentConfig `yaml:"deployments"`
}

type DeploymentConfig struct {
	MinReplicas        int32 `yaml:"minReplicas"`
	MaxReplicas        int32 `yaml:"maxReplicas"`
	ScaleUpBatchSize   int32 `yaml:"scaleUpBatchSize"`
	ScaleDownBatchSize int32 `yaml:"scaleDownBatchSize"`
}

var path = "/app/config.yaml"
var AppConfig *Config

func LoadConfigFromFile() {

	file, err := os.Open(getEnv("CONFIG_PATH", path))
	if err != nil {
		fmt.Println("Error opening config file: ", err)
		os.Exit(1)
	}
	defer file.Close()
	err = yaml.NewDecoder(file).Decode(&AppConfig)
	if err != nil {
		fmt.Println("Error decoding config file: ", err)
		os.Exit(1)
	}
	fmt.Println("Successfully loaded config from file: ", *AppConfig)
}

func getEnv(key string, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}

func (c *Config) IfNamespaceExists(namespace string) bool {
	for k := range c.Namespaces {
		if k == namespace {
			return true
		}
	}
	return false
}

func (c *Config) IfDeploymentExists(namespace string, deployment string) bool {
	for k := range c.Namespaces[namespace].Deployment {
		if k == deployment {
			return true
		}
	}
	return false
}
