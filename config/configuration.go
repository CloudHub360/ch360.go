package config

import (
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	ConfigurationRoot *ConfigurationRoot `yaml:"configuration"`
}

type ConfigurationRoot struct {
	Credentials ApiCredentialsList `yaml:"credentials"`
}

type ApiCredentialsList []ApiCredentials

type ApiCredentials struct {
	Key    string `yaml:"key"`
	Url    string `yaml:"url"`
	Id     string `yaml:"client_id"`
	Secret string `yaml:"client_secret"`
}

func NewConfiguration(clientId string, clientSecret string) *Configuration {
	var credentials = make(ApiCredentialsList, 1)

	credentials[0] = ApiCredentials{
		Key:    "default", // These credentials are the ones used by default
		Url:    "default", // These credentials are for the production API
		Id:     clientId,
		Secret: clientSecret,
	}

	config := &ConfigurationRoot{
		Credentials: credentials,
	}

	configuration := &Configuration{
		ConfigurationRoot: config,
	}

	return configuration
}

func (config *Configuration) Serialise() ([]byte, error) {
	yaml, err := yaml.Marshal(config)
	return yaml, err
}

func DeserialiseConfiguration(data []byte) (*Configuration, error) {
	var configuration Configuration
	err := yaml.Unmarshal(data, &configuration)
	return &configuration, err
}
