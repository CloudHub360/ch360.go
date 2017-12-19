package config

import (
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Credentials ApiCredentialsList `yaml:"credentials"`
}

type ApiCredentialsList []ApiCredentials

type ApiCredentials struct {
	Key    string `yaml:"key"`
	Url    string `yaml:"url"`
	Id     string `yaml:"clientId"`
	Secret string `yaml:"clientSecret"`
}

func NewConfiguration(clientId string, clientSecret string) *Configuration {
	var credentials = make(ApiCredentialsList, 1)

	credentials[0] = ApiCredentials{
		Key:    "default", // These credentials are the ones used by default
		Url:    "default", // These credentials are for the production API
		Id:     clientId,
		Secret: clientSecret,
	}

	return &Configuration{
		Credentials: credentials,
	}
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
