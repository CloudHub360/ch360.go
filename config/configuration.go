package config

import (
	"encoding/json"
)

type Configuration struct {
	ConfigurationRoot *ConfigurationRoot `json:"configuration"`
}

type ConfigurationRoot struct {
	Credentials ApiCredentialsList `json:"credentials"`
}

type ApiCredentialsList []ApiCredentials

type ApiCredentials struct {
	Key    string `json:"key"`
	Url    string `json:"url"`
	Id     string `json:"client_id"`
	Secret string `json:"client_secret"`
}

//go:generate mockery -name "FileWriter"
type FileWriter interface {
	WriteFile(filepath string, data []byte) error
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

func (config *Configuration) Save(configDir FileWriter) error {
	json, _ := json.Marshal(config)

	return configDir.WriteFile("config.json", json)
}
