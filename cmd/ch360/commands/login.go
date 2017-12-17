package commands

import "github.com/CloudHub360/ch360.go/config"

type Login struct {
	client Getter
}

func NewLogin(client Getter) *Login {
	return &Login{
		client: client,
	}
}

func (cmd *Login) Execute(clientId string, clientSecret string) error {
	configurationDirectory := config.NewConfigurationDirectory(
		config.HomeDirectoryPathGetter{},
		&config.FileSystem{})
	configuration := config.NewConfiguration(clientId, clientSecret)

	return configuration.Save(configurationDirectory)
}
