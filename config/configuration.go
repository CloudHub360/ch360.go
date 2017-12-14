package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
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

func (config *Configuration) Save() error {
	json, _ := json.Marshal(config)
	fmt.Println(string(json))

	configDirectory := configurationDirectory{}
	configDirectory.CreateIfNotExists()

	filename := filepath.Join(configDirectory.GetPath(), "config.json")

	err := ioutil.WriteFile(filename, json, 0644) //TODO: Permissions?
	return err
}
