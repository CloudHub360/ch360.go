package config

import (
	"gopkg.in/yaml.v2"
	"io"
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

//go:generate mockery -name "Writer"
type Writer interface {
	io.Writer
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

func (config *Configuration) Save(configDir io.Writer) error {
	yaml, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	_, err = configDir.Write(yaml)
	return err
}
